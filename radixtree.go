package radixtree

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
)

type Tree struct {
	root node
}

type node struct {
	label []byte
	value interface{}

	children []*node
}

var noValue = &struct{}{}

func New() *Tree {
	return &Tree{
		root: node{
			value: noValue,
		},
	}
}

func (t *Tree) String() string {
	var buf []byte
	buf = append(buf, '.')
	if t.root.hasValue() {
		buf = append(buf, fmt.Sprintf(" %+v %T", t.root.value, t.root.value)...)
	}
	buf = append(buf, '\n')

	var doPrint func(p *node, leading []byte)
	doPrint = func(p *node, leading []byte) {
		for i := range p.children {
			n := p.children[i]
			buf = append(buf, leading...)
			if i < len(p.children)-1 {
				buf = append(buf, '|')
			} else {
				buf = append(buf, '`')
			}
			buf = append(buf, "-- "...)
			buf = strconv.AppendQuote(buf, string(n.label))
			if n.hasValue() {
				buf = append(buf, fmt.Sprintf(" %+v %T", n.value, n.value)...)
			}
			buf = append(buf, '\n')
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
	return string(buf)
}

func (t *Tree) Get(key []byte) (value interface{}, exists bool) {
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

func (t *Tree) Set(key []byte, value interface{}) {
	n := &t.root
	if len(key) == 0 {
		n.value = value
		return
	}
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
					value:    noValue,
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
	switch childCount {
	case 0:
		if parent.hasValue() || parent == &t.root {
			if len(parent.children) > 1 {
				parent.children = append(parent.children[:i], parent.children[i+1:]...)
			} else {
				parent.children = nil
			}
		} else {
			parentChildCount := len(parent.children)
			if parentChildCount > 2 {
				parent.children = append(parent.children[:i], parent.children[i+1:]...)
			} else if parentChildCount == 2 {
				sibling := parent.children[1-i]
				*parent = node{
					label:    append(parent.label, sibling.label...),
					value:    sibling.value,
					children: sibling.children,
				}
			}
		}
	case 1:
		child := n.children[0]
		parent.children[i] = &node{
			label:    append(n.label, child.label...),
			value:    child.value,
			children: child.children,
		}
	default: // childCount > 1
		if !n.hasValue() {
			return false
		}
		n.value = noValue
	}
	return true
}

func (t *Tree) DeleteSubtree(prefix []byte) (deleted bool) {
	parent := &t.root
	var n *node
	var i, l int
	for len(prefix) > 0 {
		i = parent.indexForPrefix(prefix)
		if i == len(parent.children) {
			return false
		}
		n = parent.children[i]
		l = commonPrefixLength(prefix, n.label)
		if l == 0 {
			return false
		}
		if l == len(prefix) {
			break
		}
		prefix = prefix[len(n.label):]
		parent = n
	}

	parentChildCount := len(parent.children)
	if parent.hasValue() || parent == &t.root {
		if parentChildCount > 1 {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
		} else {
			parent.children = nil
		}
	} else {
		if parentChildCount > 2 {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
		} else if parentChildCount == 2 {
			sibling := parent.children[1-i]
			*parent = node{
				label:    append(parent.label, sibling.label...),
				value:    sibling.value,
				children: sibling.children,
			}
		} else {
			parent.children = nil
		}
	}
	return true
}

func (n *node) hasValue() bool {
	return n.value != noValue
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

func (n *node) String() string {
	var buf []byte
	buf = strconv.AppendQuote(buf, string(n.label))
	if n.hasValue() {
		buf = append(buf, fmt.Sprintf(" %+v %T", n.value, n.value)...)
	}
	buf = append(buf, '\n')

	var doPrint func(p *node, leading []byte)
	doPrint = func(p *node, leading []byte) {
		for i := range p.children {
			n := p.children[i]
			buf = append(buf, leading...)
			if i < len(p.children)-1 {
				buf = append(buf, '|')
			} else {
				buf = append(buf, '`')
			}
			buf = append(buf, "-- "...)
			buf = strconv.AppendQuote(buf, string(n.label))
			if n.hasValue() {
				buf = append(buf, fmt.Sprintf(" %+v %T", n.value, n.value)...)
			}
			buf = append(buf, '\n')
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
	doPrint(n, nil)
	return string(buf)
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
