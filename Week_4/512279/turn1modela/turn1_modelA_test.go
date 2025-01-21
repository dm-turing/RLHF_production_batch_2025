package main

import (
	"strconv"
	"testing"
)

func BenchmarkBubbleSort(b *testing.B) {
	for _, size := range []int{10, 1000, 10000} {
		b.Run("Size="+strconv.Itoa(size), func(b *testing.B) {
			b.StopTimer()
			data := generateRandomSlice(size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				arr := make([]int, len(data))
				copy(arr, data)
				BubbleSort(arr)
			}
		})
	}
}

func BenchmarkMergeSort(b *testing.B) {
	for _, size := range []int{10, 1000, 10000} {
		b.Run("Size="+strconv.Itoa(size), func(b *testing.B) {
			b.StopTimer()
			data := generateRandomSlice(size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				arr := make([]int, len(data))
				copy(arr, data)
				comparisons := 0
				MergeSort(arr, &comparisons)
			}
		})
	}
}

func BenchmarkHeapSort(b *testing.B) {
	for _, size := range []int{10, 1000, 10000} {
		b.Run("Size="+strconv.Itoa(size), func(b *testing.B) {
			b.StopTimer()
			data := generateRandomSlice(size)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				arr := make([]int, len(data))
				copy(arr, data)
				HeapSort(arr)
			}
		})
	}
}
