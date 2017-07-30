package radixtree_test

import (
	"testing"

	"github.com/hnakamur/radixtree"
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
			want: ". 0 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 1 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("tea"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 2 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 1 int\n" +
				"   `-- \"m\" 2 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tea"), 2)
				return t
			}(),
			want: ".\n" +
				"`-- \"tea\" 2 int\n" +
				"   `-- \"m\" 1 int\n",
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
				"   |-- \"m\" 1 int\n" +
				"   `-- \"r\" 2 int\n",
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
				"   |-- \"m\" 2 int\n" +
				"   `-- \"r\" 1 int\n",
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
				"`-- \"tea\" 2 int\n" +
				"   `-- \"m\" 1 int\n" +
				"      `-- \"work\" 3 int\n",
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
				"`-- \"tea\" 1 int\n" +
				"   `-- \"m\" 2 int\n" +
				"      |-- \"mate\" 4 int\n" +
				"      `-- \"work\" 3 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("water"), 2)
				return t
			}(),
			want: ".\n" +
				"|-- \"tea\" 1 int\n" +
				"`-- \"water\" 2 int\n",
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
			want: ". 0 int\n" +
				"|-- \"te\"\n" +
				"|  |-- \"am\" 1 int\n" +
				"|  `-- \"st\" 2 int\n" +
				"`-- \"water\" 3 int\n",
		},
	}
	for i, c := range testCases {
		got := c.tree.String()
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
				"`-- \"tea\" 1 int\n",
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
				"`-- \"water\" 2 int\n",
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
				"`-- \"tea\" 1 int\n",
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
				"`-- \"team\" 2 int\n",
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
				"   |-- \"m\" 2 int\n" +
				"   `-- \"r\" 3 int\n",
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
				"`-- \"tea\" 1 int\n" +
				"   `-- \"m\" 2 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tear"), 2)
				t.Set([]byte("test"), 3)
				return t
			}(),
			key:     []byte("test"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\"\n" +
				"   |-- \"m\" 1 int\n" +
				"   `-- \"r\" 2 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("test"), 2)
				t.Set([]byte("tester"), 3)
				return t
			}(),
			key:     []byte("test"),
			deleted: true,
			result: ".\n" +
				"`-- \"te\"\n" +
				"   |-- \"am\" 1 int\n" +
				"   `-- \"ster\" 3 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("technology"), 2)
				t.Set([]byte("test"), 3)
				return t
			}(),
			key:     []byte("test"),
			deleted: true,
			result: ".\n" +
				"`-- \"te\"\n" +
				"   |-- \"am\" 1 int\n" +
				"   `-- \"chnology\" 2 int\n",
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
				"   |-- \"m\" 1 int\n" +
				"   `-- \"r\" 2 int\n",
		},
	}
	for i, c := range testCases {
		deleted := c.tree.Delete(c.key)
		if deleted != c.deleted {
			t.Errorf("deleted unmatch, caseIndex=%d, got=%v, want=%v", i, deleted, c.deleted)
		}
		got := c.tree.String()
		if got != c.result {
			t.Errorf("result unmatch, caseIndex=%d, got=\n%s, want=\n%s", i, got, c.result)
		}
	}
}

func TestDeleteSubtree(t *testing.T) {
	testCases := []struct {
		tree    *radixtree.Tree
		prefix  []byte
		deleted bool
		result  string
	}{
		{
			tree: func() *radixtree.Tree {
				return radixtree.New()
			}(),
			prefix:  []byte("tea"),
			deleted: false,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			prefix:  []byte("water"),
			deleted: false,
			result: ".\n" +
				"`-- \"tea\" 1 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			prefix:  []byte("tea"),
			deleted: true,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				return t
			}(),
			prefix:  []byte("te"),
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
			prefix:  []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"water\" 2 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			prefix:  []byte("team"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" 1 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				return t
			}(),
			prefix:  []byte("tea"),
			deleted: true,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				return t
			}(),
			prefix:  []byte("tea"),
			deleted: true,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				return t
			}(),
			prefix:  []byte("tear"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" 1 int\n" +
				"   `-- \"m\" 2 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("team"), 1)
				t.Set([]byte("tear"), 2)
				return t
			}(),
			prefix:  []byte("tea"),
			deleted: true,
			result:  ".\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				t.Set([]byte("teamwork"), 4)
				return t
			}(),
			prefix:  []byte("team"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" 1 int\n" +
				"   `-- \"r\" 3 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				t.Set([]byte("teamwork"), 4)
				t.Set([]byte("test"), 5)
				return t
			}(),
			prefix:  []byte("tea"),
			deleted: true,
			result: ".\n" +
				"`-- \"test\" 5 int\n",
		},
		{
			tree: func() *radixtree.Tree {
				t := radixtree.New()
				t.Set([]byte("tea"), 1)
				t.Set([]byte("team"), 2)
				t.Set([]byte("tear"), 3)
				t.Set([]byte("teamwork"), 4)
				t.Set([]byte("test"), 5)
				return t
			}(),
			prefix:  []byte("test"),
			deleted: true,
			result: ".\n" +
				"`-- \"tea\" 1 int\n" +
				"   |-- \"m\" 2 int\n" +
				"   |  `-- \"work\" 4 int\n" +
				"   `-- \"r\" 3 int\n",
		},
	}
	for i, c := range testCases {
		deleted := c.tree.DeleteSubtree(c.prefix)
		if deleted != c.deleted {
			t.Errorf("deleted unmatch, caseIndex=%d, got=%v, want=%v", i, deleted, c.deleted)
		}
		got := c.tree.String()
		if got != c.result {
			t.Errorf("result unmatch, caseIndex=%d, got=\n%s, want=\n%s", i, got, c.result)
		}
	}
}
