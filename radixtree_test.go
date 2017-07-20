package radixtree

import (
	"os"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	tr := Tree{root: node{
		value: []byte{'0'},
		children: []*node{
			&node{
				label: []byte("te"),
				children: []*node{
					&node{
						label: []byte("am"),
						value: []byte{'1'},
					},
					&node{
						label: []byte("st"),
						value: []byte{'2'},
					},
				},
			},
			&node{
				label: []byte("water"),
				value: []byte{'3'},
			},
		},
	}}
	tr.PrettyPrint(os.Stdout)
}
