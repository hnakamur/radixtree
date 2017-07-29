package radixtree_test

import (
	"bytes"
	"testing"

	"bitbucket.org/hnakamur/radixtree"
)

func TestPrettyPrint(t *testing.T) {
	testCases := []struct {
		tree func() *radixtree.Tree
		want string
	}{
		{
			tree: func() *radixtree.Tree {
				return radixtree.New()
			},
			want: ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte{}, 0)
				return t
			},
			want: ". 0\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			},
			want: ".\n" +
				"`-- \"tea\" 1\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte{}, 0)
				t.Set([]byte("team"), 1)
				t.Set([]byte("test"), 2)
				t.Set([]byte("water"), 3)
				return t
			},
			want: ". 0\n" +
				"|-- \"te\"\n" +
				"|  |-- \"am\" 1\n" +
				"|  `-- \"st\" 2\n" +
				"`-- \"water\" 3\n",
		},
	}
	for i, c := range testCases {
		var buf bytes.Buffer
		c.tree().PrettyPrint(&buf)
		got := buf.String()
		if got != c.want {
			t.Errorf("unmatch result, caseIndex=%d, got=\n%s, want=\n%s", i, got, c.want)
		}
	}
}
