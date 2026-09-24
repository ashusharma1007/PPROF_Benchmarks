package datastructures

import "testing"

const collectionSize = 10000

func BenchmarkBuildSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildSlice(collectionSize)
	}
}

func BenchmarkBuildMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildMap(collectionSize)
	}
}

func BenchmarkBuildLinkedList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BuildLinkedList(collectionSize)
	}
}

// The search benchmarks target the last element: worst case for a linear scan.

func BenchmarkSliceSearch(b *testing.B) {
	s := BuildSlice(collectionSize)
	target := collectionSize - 1
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SliceSearch(s, target)
	}
}

func BenchmarkMapSearch(b *testing.B) {
	m := BuildMap(collectionSize)
	target := collectionSize - 1
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MapSearch(m, target)
	}
}

func BenchmarkLinkedListSearch(b *testing.B) {
	head := BuildLinkedList(collectionSize)
	target := collectionSize - 1
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		LinkedListSearch(head, target)
	}
}
