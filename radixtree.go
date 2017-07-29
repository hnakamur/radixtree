package radixtree

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
)

type Tree struct {
	root node
}

type node struct {
	label []byte
	// To save memory, we don't a flag field like hasValue.
	// If this node does not have value, value is nil.
	// So you cannot use nil as a value.
	// Note you can use an empty byte slice instead.
	value []byte

	children []*node
}

func (n *node) hasValue() bool {
	return n.value != nil
}

func (t Tree) PrettyPrint(w io.Writer) {
	buf := make([]byte, 0, 80)
	buf = append(buf[:0], '.')
	if t.root.hasValue() {
		buf = append(buf, fmt.Sprintf(" %q", string(t.root.value))...)
	}
	buf = append(buf, '\n')
	w.Write(buf)

	var doPrint func(p *node, leading []byte)
	doPrint = func(p *node, leading []byte) {
		for i := range p.children {
			n := p.children[i]
			buf = append(buf[:0], leading...)
			if i < len(p.children)-1 {
				buf = append(buf, '|')
			} else {
				buf = append(buf, '`')
			}
			buf = append(buf, "-- "...)
			buf = strconv.AppendQuote(buf, string(n.label))
			if n.hasValue() {
				buf = append(buf, ' ')
				buf = strconv.AppendQuote(buf, string(n.value))
			}
			buf = append(buf, '\n')
			w.Write(buf)
			if len(n.children) > 0 {
				var leading2 []byte
				if i < len(p.children)-1 {
					leading2 = append(leading, "|  "...)
				} else {
					leading2 = append(leading, "   "...)
				}
				doPrint(n, leading2)
			}
		}
	}
	doPrint(&t.root, nil)
}

// edge represents a child node.
// childIndex is the index of the child in parent's child nodes.
type edge struct {
	parent     *node
	childIndex int
}

func (n edge) child() *node {
	return n.parent.children[n.childIndex]
}

func (n *node) indexForPrefix(prefix []byte) int {
	f := func(i int) bool {
		label := n.children[i].label
		if commonPrefixLength(label, prefix) > 0 {
			return true
		}
		return bytes.Compare(label, prefix) >= 0
	}
	return sort.Search(len(n.children), f)
}

type path struct {
	tree  *Tree
	edges []edge
	node  *node
}

func (p path) depth() int {
	return len(p.edges)
}

func (p path) nodeAtDepth(depth int) *node {
	if depth == 0 {
		return &p.tree.root
	}
	if depth == len(p.edges)-1 {
		return p.node
	}
	return p.edges[depth-1].child()
}

func (t *Tree) Get(key []byte) (value []byte, exists bool) {
	prefix := key
	n := &t.root
	for len(prefix) > 0 {
		i := n.indexForPrefix(prefix)
		if i == len(n.children) || !bytes.HasPrefix(prefix, n.children[i].label) {
			return nil, false
		}
		prefix = prefix[len(n.children[i].label):]
		n = n.children[i]
	}
	return n.value, true
}

func (t *Tree) Set(key, value []byte) {
	n := &t.root
	prefix := key
	for len(prefix) > 0 {
		i := n.indexForPrefix(prefix)
		if i == len(n.children) {
			n.children = append(n.children, &node{label: prefix, value: value})
			return
		}
		child := n.children[i]
		childLabel := child.label
		l := commonPrefixLength(prefix, childLabel)
		if l == 0 {
			// Insert new node at i'th children
			n.children = append(n.children, nil)
			copy(n.children[i+1:], n.children[i:])
			n.children[i] = &node{label: prefix, value: value}
			return
		}
		if l < len(prefix) {
			if l < len(childLabel) {
				myRestLabel := prefix[l:]
				child.label = childLabel[l:]
				var children []*node
				if bytes.Compare(myRestLabel, child.label) < 0 {
					children = []*node{&node{label: myRestLabel, value: value}, child}
				} else {
					children = []*node{child, &node{label: myRestLabel, value: value}}
				}
				n.children[i] = &node{
					label:    prefix[:l],
					children: children,
				}
				return
			}
		} else { // l == len(prefix)
			if l < len(childLabel) {
				child.label = childLabel[l:]
				n.children[i] = &node{
					label:    prefix,
					value:    value,
					children: []*node{child},
				}
			} else { // l == len(childLabel)
				n.children[i].value = value
			}
			return
		}
		prefix = prefix[len(childLabel):]
		n = child
	}
}

func (t *Tree) Delete(key []byte) (deleted bool) {
	parent := &t.root
	prefix := key
	var n *node
	var i int
	for len(prefix) > 0 {
		i = parent.indexForPrefix(prefix)
		if i == len(parent.children) {
			return false
		}
		n = parent.children[i]
		l := commonPrefixLength(prefix, n.label)
		if l == 0 {
			return false
		}
		if l == len(prefix) {
			break
		}
		prefix = prefix[len(n.label):]
		parent = n
	}

	childCount := len(n.children)
	if childCount == 0 {
		if len(parent.children) > 1 {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
		} else {
			parent.children = nil
		}
	} else if childCount == 1 {
		child := n.children[0]
		child.label = append(n.label, child.label...)
		parent.children[i] = child
	} else { // childCount > 1
		n.value = nil
	}
	return true
}

func (t *Tree) pathForPrefix(prefix []byte) path {
	n := &t.root
	p := path{tree: t, node: n}
	for len(prefix) > 0 {
		i := n.indexForPrefix(prefix)
		p.edges = append(p.edges, edge{parent: n, childIndex: i})
		if i < len(n.children) && bytes.HasPrefix(prefix, n.children[i].label) {
			p.node = n.children[i]
		} else {
			p.node = nil
			break
		}
		prefix = prefix[len(n.children[i].label):]
		n = n.children[i]
	}
	return p
}

func commonPrefixLength(a, b []byte) int {
	prefix := a
	other := b
	if len(b) < len(a) {
		prefix = b
		other = a
	}
	for i, c := range prefix {
		if c != other[i] {
			return i
		}
	}
	return len(prefix)
}
