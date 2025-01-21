
package main  
import (  
    "fmt"
    "math/rand"
    "runtime"
    "testing"
)
// these are global variables to be used in the benchmark functions
var smallArray []int
var mediumArray []int
var largeArray []int
// Function to generate a random array of size n
func generateRandomArray(n int) []int {
    array := make([]int, n)
    for i := 0; i < n; i++ {
        array[i] = rand.Intn(1000)
    }
    return array
}
// Function to initialize the global arrays
func init() {
    smallArray = generateRandomArray(10000)
    mediumArray = generateRandomArray(100000)
    largeArray = generateRandomArray(1000000)
}
// Function to measure the execution time of a function
func measureExecutionTime(f func()) time.Duration {
    start := time.Now()
    f()
    return time.Since(start)
}
// Function to measure the memory usage of a function
func measureMemoryUsage(f func()) int64 {
    var memStats runtime.MemStats
    runtime.ReadMemStats(&memStats)
    start := memStats.Alloc
    f()
    runtime.ReadMemStats(&memStats)
    return memStats.Alloc - start
}
// BenchmarkBubbleSort-8 measures the execution time, memory usage, number of comparisons, and swaps for bubble sort.
func BenchmarkBubbleSortSmall(b *testing.B) {
    array := smallArray
    for i := 0; i < b.N; i++ {
        BubbleSort(array)
    }
}
func BenchmarkBubbleSortMedium(b *testing.B) {
    array := mediumArray
    for i := 0; i < b.N; i++ {   BubbleSort(array)
    }
}
func BenchmarkBubbleSortLarge(b *testing.B) {
    array := largeArray   for i := 0; i < b.N; i++ {
        BubbleSort(array)
    }
}
// BenchmarkMergeSort-8 measures the execution time, memory usage, number of comparisons, and swaps for merge sort.
func BenchmarkMergeSortSmall(b *testing.B) {
    array := smallArray
    for i := 0; i < b.N; i++ { 
        MergeSort(array)
    }
}
func BenchmarkMergeSortMedium(b *testing.B) {
    array := mediumArray
    for i := 0; i < b.N; i++ {
        MergeSort(array)
    }
}
func BenchmarkMergeSortLarge(b *testing.B) {
    array := largeArray
    for i := 0; i < b.N; i++ {
        MergeSort(array)
    }
}
// BenchmarkHeapSort-8 measures the execution time, memory usage, number of comparisons, and swaps for heap sort.
func BenchmarkHeapSortSmall(b *testing.B) {  
    array := smallArray
    for i := 0; i < b.N; i++ {
        HeapSort(array)
    }
}
func BenchmarkHeapSortMedium(b *testing.B) {
    array := mediumArray
    for i := 0; i < b.N; i++ {
        HeapSort(array)
    }
}
func BenchmarkHeapSortLarge(b *testing.B) {  
    array := largeArray
    for i := 0; i < b.N; i++ {
        HeapSort(array)
    }
}
// Bubble sort implementation
func BubbleSort(array []int) []int {
    n := len(array)
    for i := 0; i < n-1; i++ {
        swapped := false
        for j := 0; j < n-i-1; j++ {
            if array[j] > array[j+1] {
                array[j], array[j+1] = array[j+1], array[j]
                swapped = true
            }