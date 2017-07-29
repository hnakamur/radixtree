package radixtree

import (
	"log"
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

	doGet := func(key string) {
		val, exist := tr.Get([]byte(key))
		log.Printf("get key=%q, val=%q, exist=%v", key, string(val), exist)
	}
	//doGet("")
	doGet("tea")
	//doGet("team")
	//doGet("test")
	//doGet("tess")
	//doGet("testable")
	//doGet("water")
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
