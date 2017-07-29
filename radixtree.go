package radixtree

import (
	"bytes"
	"fmt"
	"io"
	"log"
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
	p := t.pathForPrefix(key)
	if p.node == nil {
		return nil, false
	}
	return (*p.node).value, true
}

func (t *Tree) pathForPrefix(prefix []byte) path {
	n := &t.root
	p := path{tree: t, node: n}
	for len(prefix) > 0 {
		f := func(i int) bool {
			log.Printf("i=%d, childLabel=%q, prefix=%q, fRes=%v", i, string(n.children[i].label), string(prefix), bytes.Compare(n.children[i].label, prefix) >= 0)
			return bytes.Compare(n.children[i].label, prefix) >= 0
		}
		i := sort.Search(len(n.children), f)
		if i > 0 && commonPrefixLength(prefix, n.children[i-1].label) > 0 {
			log.Printf("decrement i since it has common prefix with previous chlidren")
			i--
		}
		log.Printf("search result, prefix=%q, i=%d", prefix, i)
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
