package memalloc

import "testing"

func BenchmarkConcatWithPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatWithPlus(10000)
	}
}

func BenchmarkConcatWithBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatWithBuilder(10000)
	}
}

func BenchmarkConcatWithBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ConcatWithBytes(10000)
	}
}
