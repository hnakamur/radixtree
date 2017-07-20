package radixtree

import (
	"fmt"
	"io"
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
