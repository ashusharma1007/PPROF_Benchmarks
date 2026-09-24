package datastructures

// BuildSlice returns [0, 1, ..., n-1] in one contiguous allocation.
func BuildSlice(n int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = i
	}
	return s
}

// BuildMap returns {0: true, ..., n-1: true}. Costs more memory and more
// allocations than a slice, in exchange for O(1) lookups.
func BuildMap(n int) map[int]bool {
	m := make(map[int]bool, n)
	for i := 0; i < n; i++ {
		m[i] = true
	}
	return m
}

// SliceSearch scans linearly. O(n), but contiguous memory keeps it cache-friendly.
func SliceSearch(s []int, target int) bool {
	for _, v := range s {
		if v == target {
			return true
		}
	}
	return false
}

// MapSearch is a single hash lookup. O(1).
func MapSearch(m map[int]bool, target int) bool {
	return m[target]
}

// Node is 16 bytes: an int plus a pointer. Each one is a separate heap
// allocation, so building a list of n nodes costs n allocations.
type Node struct {
	Value int
	Next  *Node
}

// BuildLinkedList creates a list of n nodes.
func BuildLinkedList(n int) *Node {
	var head *Node
	for i := n - 1; i >= 0; i-- {
		head = &Node{Value: i, Next: head}
	}
	return head
}

// LinkedListSearch walks node by node. Also O(n), but ~2x slower than
// SliceSearch because scattered nodes defeat the CPU prefetcher.
func LinkedListSearch(head *Node, target int) bool {
	current := head
	for current != nil {
		if current.Value == target {
			return true
		}
		current = current.Next
	}
	return false
}
