package sorting

// BubbleSort sorts a copy of arr in O(n²) time, comparing every adjacent pair.
func BubbleSort(arr []int) []int {
	n := len(arr)
	result := make([]int, n)
	copy(result, arr)

	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	return result
}

// QuickSort sorts a copy of arr in O(n log n) time by divide and conquer.
func QuickSort(arr []int) []int {
	n := len(arr)
	result := make([]int, n)
	copy(result, arr)
	quickSortHelper(result, 0, n-1)
	return result
}

func quickSortHelper(arr []int, low, high int) {
	if low < high {
		pivotIndex := partition(arr, low, high)
		quickSortHelper(arr, low, pivotIndex-1)
		quickSortHelper(arr, pivotIndex+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1

	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
