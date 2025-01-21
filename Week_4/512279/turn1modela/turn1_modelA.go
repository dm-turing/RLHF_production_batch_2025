package main

import (
	"math/rand"
	"time"
)

// Bubble sort implementation
func BubbleSort(arr []int) (int, int) {
	n := len(arr)
	comparisons := 0
	swaps := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			comparisons++
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swaps++
			}
		}
	}
	return comparisons, swaps
}

// Merge sort implementation
func MergeSort(arr []int, comparisons *int) []int {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	left := MergeSort(arr[:mid], comparisons)
	right := MergeSort(arr[mid:], comparisons)

	return merge(left, right, comparisons)
}

func merge(left, right []int, comparisons *int) []int {
	result := make([]int, 0)
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		*comparisons++
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// Heap sort implementation
func HeapSort(arr []int) (int, int) {
	n := len(arr)
	swaps := 0
	comparisons := 0
	for i := n/2 - 1; i >= 0; i-- {
		comparisons, swaps = heapify(arr, n, i, comparisons, swaps)
	}
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		swaps++
		comparisons, swaps = heapify(arr, i, 0, comparisons, swaps)
	}
	return comparisons, swaps
}

func heapify(arr []int, n int, i int, comparisons int, swaps int) (int, int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n {
		comparisons++
		if arr[left] > arr[largest] {
			largest = left
		}
	}

	if right < n {
		comparisons++
		if arr[right] > arr[largest] {
			largest = right
		}
	}

	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		swaps++
		comparisons, swaps = heapify(arr, n, largest, comparisons, swaps)
	}

	return comparisons, swaps
}

func generateRandomSlice(size int) []int {
	rand.Seed(time.Now().UnixNano())
	slice := make([]int, size)
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(10000) // generate random numbers between 0 and 10000
	}
	return slice
}

func main() {} // Making this file a standalone package
