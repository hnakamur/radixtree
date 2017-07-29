package radixtree

import (
	"bytes"
	"testing"
)

func TestGet(t *testing.T) {
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

	testCases := []struct {
		key    []byte
		value  []byte
		exists bool
	}{
		{[]byte(""), []byte{'0'}, true},
		{[]byte("tea"), nil, false},
		{[]byte("team"), []byte{'1'}, true},
		{[]byte("tear"), nil, false},
		{[]byte("test"), []byte{'2'}, true},
		{[]byte("testable"), nil, false},
		{[]byte("water"), []byte{'3'}, true},
	}
	for _, c := range testCases {
		value, exists := tr.Get(c.key)
		if exists != c.exists {
			t.Errorf("exists unmatch, got=%v, want=%v", exists, c.exists)
		} else if !bytes.Equal(value, c.value) {
			t.Errorf("value unmatch, got=%s, want=%s", string(value), string(c.value))
		}
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		tree   Tree
		key    []byte
		value  []byte
		result string
	}{
		{
			tree:  Tree{},
			key:   []byte("tea"),
			value: []byte("1"),
			result: ".\n" +
				"`-- \"tea\" \"1\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("coffee"),
			value: []byte("2"),
			result: ".\n" +
				"|-- \"coffee\" \"2\"\n" +
				"`-- \"tea\" \"1\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("water"),
			value: []byte("2"),
			result: ".\n" +
				"|-- \"tea\" \"1\"\n" +
				"`-- \"water\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("tea"),
			value: []byte("2"),
			result: ".\n" +
				"`-- \"tea\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("team"),
			value: []byte("2"),
			result: ".\n" +
				"`-- \"tea\" \"1\"\n" +
				"   `-- \"m\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("team"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("tea"),
			value: []byte("2"),
			result: ".\n" +
				"`-- \"tea\" \"2\"\n" +
				"   `-- \"m\" \"1\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("team"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("tear"),
			value: []byte("2"),
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" \"1\"\n" +
				"   `-- \"r\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tear"),
							value: []byte("1"),
						},
					},
				},
			},
			key:   []byte("team"),
			value: []byte("2"),
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" \"2\"\n" +
				"   `-- \"r\" \"1\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("team"),
							value: []byte("1"),
						},
						&node{
							label: []byte("water"),
							value: []byte("2"),
						},
					},
				},
			},
			key:   []byte("tear"),
			value: []byte("3"),
			result: ".\n" +
				"|-- \"tea\"\n" +
				"|  |-- \"m\" \"1\"\n" +
				"|  `-- \"r\" \"3\"\n" +
				"`-- \"water\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte("1"),
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte("2"),
								},
							},
						},
					},
				},
			},
			key:   []byte("teamwork"),
			value: []byte("3"),
			result: ".\n" +
				"`-- \"tea\" \"1\"\n" +
				"   `-- \"m\" \"2\"\n" +
				"      `-- \"work\" \"3\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte("1"),
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte("2"),
									children: []*node{
										&node{
											label: []byte("work"),
											value: []byte("3"),
										},
									},
								},
							},
						},
					},
				},
			},
			key:   []byte("teammate"),
			value: []byte("4"),
			result: ".\n" +
				"`-- \"tea\" \"1\"\n" +
				"   `-- \"m\" \"2\"\n" +
				"      |-- \"mate\" \"4\"\n" +
				"      `-- \"work\" \"3\"\n",
		},
	}
	for i, c := range testCases {
		c.tree.Set(c.key, c.value)
		var buf bytes.Buffer
		c.tree.PrettyPrint(&buf)
		got := buf.String()
		if got != c.result {
			t.Errorf("result unmatch, caseIndex=%d, got=\n%s, want=\n%s", i, got, c.result)
		}
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		tree    Tree
		key     []byte
		deleted bool
		result  string
	}{
		{
			tree:    Tree{},
			key:     []byte("tea"),
			deleted: false,
			result:  ".\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
						},
					},
				},
			},
			key:     []byte("water"),
			deleted: false,
			result: ".\n" +
				"`-- \"tea\" \"1\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
						},
					},
				},
			},
			key:     []byte("tea"),
			deleted: true,
			result:  ".\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
						},
						&node{
							label: []byte("water"),
							value: []byte{'2'},
						},
					},
				},
			},
			key:     []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"water\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte{'2'},
								},
							},
						},
					},
				},
			},
			key:     []byte("team"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" \"1\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte{'2'},
								},
							},
						},
					},
				},
			},
			key:     []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"team\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte{'2'},
								},
								&node{
									label: []byte("r"),
									value: []byte{'3'},
								},
							},
						},
					},
				},
			},
			key:     []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" \"2\"\n" +
				"   `-- \"r\" \"3\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte{'2'},
								},
								&node{
									label: []byte("r"),
									value: []byte{'3'},
								},
							},
						},
					},
				},
			},
			key:     []byte("tear"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" \"1\"\n" +
				"   `-- \"m\" \"2\"\n",
		},
		{
			tree: Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							children: []*node{
								&node{
									label: []byte("m"),
									value: []byte{'2'},
								},
								&node{
									label: []byte("r"),
									value: []byte{'3'},
								},
							},
						},
					},
				},
			},
			key:     []byte("tea"),
			deleted: false,
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" \"2\"\n" +
				"   `-- \"r\" \"3\"\n",
		},
	}
	for i, c := range testCases {
		deleted := c.tree.Delete(c.key)
		if deleted != c.deleted {
			t.Errorf("deleted unmatch, caseIndex=%d, got=%v, want=%v", i, deleted, c.deleted)
		}
		var buf bytes.Buffer
		c.tree.PrettyPrint(&buf)
		got := buf.String()
		if got != c.result {
			t.Errorf("result unmatch, caseIndex=%d, got=\n%s, want=\n%s", i, got, c.result)
		}
	}
}

func TestPrettyPrint(t *testing.T) {
	testCases := []struct {
		tree Tree
		want string
	}{
		{
			Tree{},
			".\n",
		},
		{
			Tree{
				root: node{
					value: []byte{'0'},
				},
			},
			". \"0\"\n",
		},
		{
			Tree{
				root: node{
					children: []*node{
						&node{
							label: []byte("tea"),
							value: []byte{'1'},
						},
					},
				},
			},
			".\n" +
				"`-- \"tea\" \"1\"\n",
		},
		{
			Tree{
				root: node{
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
				},
			},
			". \"0\"\n" +
				"|-- \"te\"\n" +
				"|  |-- \"am\" \"1\"\n" +
				"|  `-- \"st\" \"2\"\n" +
				"`-- \"water\" \"3\"\n",
		},
	}
	for i, c := range testCases {
		var buf bytes.Buffer
		c.tree.PrettyPrint(&buf)
		got := buf.String()
		if got != c.want {
			t.Errorf("unmatch result, caseIndex=%d, got=\n%s, want=\n%s", i, got, c.want)
		}
	}
}

func TestCommonPrefixLength(t *testing.T) {
	testCases := []struct {
		a, b []byte
		want int
	}{
		{nil, nil, 0},
		{[]byte{}, []byte{}, 0},
		{[]byte("te"), []byte("tea"), 2},
		{[]byte("tea"), []byte("te"), 2},
		{[]byte("test"), []byte("toast"), 1},
		{[]byte("team"), []byte("test"), 2},
		{[]byte("test"), []byte("water"), 0},
	}
	for _, tc := range testCases {
		got := commonPrefixLength(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("a=%q, b=%q, got=%d, want=%d", tc.a, tc.b, got, tc.want)
		}
	}
}
