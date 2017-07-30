package radixtree

import (
	"bytes"
	"testing"
)

func TestPrettyPrint(t *testing.T) {
	testCases := []struct {
		tree Tree
		want string
	}{
		{
			Tree{
				root: node{
					value: noValue,
				},
			},
			".\n",
		},
		{
			Tree{
				root: node{
					value: "0",
				},
			},
			". 0 string\n",
		},
		{
			Tree{
				root: node{
					value: noValue,
					children: []*node{
						&node{
							label: []byte("tea"),
							value: "1",
						},
					},
				},
			},
			".\n" +
				"`-- \"tea\" 1 string\n",
		},
		{
			Tree{
				root: node{
					value: 0,
					children: []*node{
						&node{
							label: []byte("te"),
							value: noValue,
							children: []*node{
								&node{
									label: []byte("am"),
									value: 1,
								},
								&node{
									label: []byte("st"),
									value: 2,
								},
							},
						},
						&node{
							label: []byte("water"),
							value: 3,
						},
					},
				},
			},
			". 0 int\n" +
				"|-- \"te\"\n" +
				"|  |-- \"am\" 1 int\n" +
				"|  `-- \"st\" 2 int\n" +
				"`-- \"water\" 3 int\n",
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
