package itools

import (
	"iter"
)

func Count[T any](i iter.Seq[T]) uint64 {
	var cnt uint64 = 0
	for range i {
		cnt += 1
	}
	return cnt
}
