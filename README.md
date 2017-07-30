radixtree [![Build Status](https://travis-ci.org/hnakamur/radixtree.png)](https://travis-ci.org/hnakamur/radixtree) [![Go Report Card](https://goreportcard.com/badge/github.com/hnakamur/radixtree)](https://goreportcard.com/report/github.com/hnakamur/radixtree) [![GoDoc](https://godoc.org/github.com/hnakamur/radixtree?status.svg)](https://godoc.org/github.com/hnakamur/radixtree)
=========

Package radixtree provides a simple and straightforward implementation
for radixtree.

This implementation uses the binary search to find the index in children
at each level of nodes in a tree. So it will be slower than map when
the level becomes large.

The advantage of the radixtree implementation is the cost of deleting
a subtree for a prefix is cheap. It is roughly same as deleting a single
key.

This implementation is not goroutine safe, so you need to use a lock in your
code when multiple goroutines concurrently access the same tree.
