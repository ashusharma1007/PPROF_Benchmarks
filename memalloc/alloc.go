package memalloc

import "strings"

// ConcatWithPlus builds a string with +=. Strings are immutable, so every
// iteration allocates a new string and copies the old one: n allocations and
// O(n²) bytes copied.
func ConcatWithPlus(n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += "a"
	}
	return result
}

// ConcatWithBuilder uses strings.Builder with a pre-grown buffer: 1 allocation.
func ConcatWithBuilder(n int) string {
	var builder strings.Builder
	builder.Grow(n)
	for i := 0; i < n; i++ {
		builder.WriteByte('a')
	}
	return builder.String()
}

// ConcatWithBytes fills a pre-sized []byte and converts once. Same single
// allocation as Builder, with less per-byte bookkeeping.
func ConcatWithBytes(n int) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = 'a'
	}
	return string(buf)
}
