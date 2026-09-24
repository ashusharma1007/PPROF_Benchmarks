package sorting

import (
	"math/rand"
	"testing"
)

func generateSlice(n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rand.Intn(10000)
	}
	return arr
}

func BenchmarkBubbleSort(b *testing.B) {
	arr := generateSlice(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BubbleSort(arr)
	}
}

func BenchmarkQuickSort(b *testing.B) {
	arr := generateSlice(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(arr)
	}
}
