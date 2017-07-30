package radixtree_test

import (
	"bytes"
	"testing"

	"bitbucket.org/hnakamur/radixtree"
)

func TestSet(t *testing.T) {
	testCases := []struct {
		tree *radixtree.Tree
		want string
	}{
		{
			tree: func() *radixtree.Tree {
				return radixtree.New()
			}(),
			want: ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte{}, 0)
				return t
			}(),
			want: ". 0\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 1\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("tea"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 1\n" +
				"   `-- \"m\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tea"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 2\n" +
				"   `-- \"m\" 1\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tear"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" 1\n" +
				"   `-- \"r\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tear"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" 2\n" +
				"   `-- \"r\" 1\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tea"), 2)
				t.Set([]byte("teamwork"), 3)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 2\n" +
				"   `-- \"m\" 1\n" +
				"      `-- \"work\" 3\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("teamwork"), 3)
				t.Set([]byte("teammate"), 4)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 1\n" +
				"   `-- \"m\" 2\n" +
				"      |-- \"mate\" 4\n" +
				"      `-- \"work\" 3\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("water"), 2)
				return t
			}(),
			want: ".\n" +
				"|-- \"tea\" 1\n" +
				"`-- \"water\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte{}, 0)
				t.Set([]byte("team"), 1)
				t.Set([]byte("test"), 2)
				t.Set([]byte("water"), 3)
				return t
			}(),
			want: ". 0\n" +
				"|-- \"te\"\n" +
				"|  |-- \"am\" 1\n" +
				"|  `-- \"st\" 2\n" +
				"`-- \"water\" 3\n",
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

func TestGet(t *testing.T) {
	tree1 := func() *radixtree.Tree {
		t := radixtree.New()
		t.Set([]byte{}, 0)
		t.Set([]byte("team"), 1)
		t.Set([]byte("test"), 2)
		t.Set([]byte("water"), 3)
		return t
	}()
	testCases := []struct {
		tree   *radixtree.Tree
		key    []byte
		value  int
		exists bool
	}{
		{tree: tree1, key: []byte{}, value: 0, exists: true},
		{tree: tree1, key: []byte("team"), value: 1, exists: true},
		{tree: tree1, key: []byte("test"), value: 2, exists: true},
		{tree: tree1, key: []byte("water"), value: 3, exists: true},
		{tree: tree1, key: []byte("tea"), exists: false},
	}
	for i, c := range testCases {
		value, exists := c.tree.Get(c.key)
		if exists != c.exists {
			t.Errorf("exists unmatch, caseIndex=%d, got=%v, want=%v", i, exists, c.exists)
		}
		if exists {
			intVal, ok := value.(int)
			if !ok {
				t.Errorf("non int value, caseIndex=%d", i)
			} else if intVal != c.value {
				t.Errorf("value unmatch, caseIndex=%d, got=%d, want=%d", i, intVal, c.value)
			}
		}
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		tree    *radixtree.Tree
		key     []byte
		deleted bool
		result  string
	}{
		{
			tree: func() *radixtree.Tree {
				return radixtree.New()
			}(),
			key:     []byte("tea"),
			deleted: false,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			key:     []byte("water"),
			deleted: false,
			result: ".\n" +
				"`-- \"tea\" 1\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			key:     []byte("tea"),
			deleted: true,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("water"), 2)
				return t
			}(),
			key:     []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"water\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			key:     []byte("team"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" 1\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			key:     []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"team\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				return t
			}(),
			key:     []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" 2\n" +
				"   `-- \"r\" 3\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				return t
			}(),
			key:     []byte("tear"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" 1\n" +
				"   `-- \"m\" 2\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tear"), 2)
				return t
			}(),
			key:     []byte("tea"),
			deleted: false,
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" 1\n" +
				"   `-- \"r\" 2\n",
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
