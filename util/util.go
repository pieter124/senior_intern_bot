package util

type Empty struct{}
type Set[T comparable] map[T]Empty
