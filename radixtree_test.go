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
	var buf bytes.Buffer
	tr.PrettyPrint(&buf)
	got := buf.String()
	want := ". \"0\"\n" +
		"|-- \"te\"\n" +
		"|  |-- \"am\" \"1\"\n" +
		"|  `-- \"st\" \"2\"\n" +
		"`-- \"water\" \"3\"\n"
	if got != want {
		t.Errorf("unmatch result, got=\n%s, want=\n%s", got, want)
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
