
package main  
import (  
    "fmt"
    "math/rand"
    "sort"
    "time"
)

// Use bubble sort to sort an array of integers
func bubbleSort(arr []int) {  
    n := len(arr)  
    for i := 0; i < n-1; i++ {  
        for j := 0; j < n-i-1; j++ {  
            if arr[j] > arr[j+1] {  
                arr[j], arr[j+1] = arr[j+1], arr[j]  
            }  
        }  
    }  
}  


// Use merge sort to sort an array of integers
func mergeSort(arr []int) []int {  
    if len(arr) <= 1 {  
        return arr  
    }  
    mid := len(arr) / 2  
    left := mergeSort(arr[:mid])  
    right := mergeSort(arr[mid:])  
    return merge(left, right)  
}  


func merge(left, right []int) []int {  
    result := make([]int, 0, len(left)+len(right))  
    i, j := 0, 0  
    for i < len(left) && j < len(right) {  
        if left[i] < right[j] {  
            result = append(result, left[i])  
            i++  
        } else {  
            result = append(result, right[j])  
            j++  
        }  
    }  
    for i < len(left) {  
        result = append(result, left[i])  
        i++  
    }  
    for j < len(right) {  
        result = append(result, right[j])  
        j++  
    }  
    return result  
}  


// Use the built-in sort.Slice function to sort an array of integers
func sortSlice(arr []int) {  
    sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })  
}  

func main() {  
    rand.Seed(time.Now().UnixNano())  
    const size = 100000  
    originalArray := make([]int, size)  
    for i := range originalArray {  
        originalArray[i] = rand.Intn(size)  
    }  
    
    //Analysis of Bubble Sort
    fmt.Println("Analysis of Bubble Sort:")
    start := time.Now()  
    bubbleSortArray := make([]int, len(originalArray))  
    copy(bubbleSortArray, originalArray)  
    bubbleSort(bubbleSortArray)  
    bubbleSortTime := time.Since(start)  
    fmt.Printf("Time taken: %s\n", bubbleSortTime)  
    isSorted := sort.IntsAreSorted(bubbleSortArray)  
    fmt.Println("Is Sorted:", isSorted)  

    //Analysis of Merge Sort
    fmt.Println("\nAnalysis of Merge Sort:")  
    start = time.Now()  
    mergeSortArray := make([]int, len(originalArray))  
    copy(mergeSortArray, originalArray)  
    mergeSortArray = mergeSort(mergeSortArray)  
    mergeSortTime := time.Since(start)  
    fmt.Printf("Time taken: %s\n", mergeSortTime)  
    isSorted = sort.IntsAreSorted(mergeSortArray)  
    fmt.Println("Is Sorted:", isSorted)  

    //Analysis of sort.Slice function
    fmt.Println("\nAnalysis of sort.Slice function:")  
    start = time.Now()  
    sortSliceArray := make([]int, len(originalArray))  
    copy(sortSliceArray, originalArray)  
    sortSlice(sortSliceArray)  
    sortSliceTime := time.Since(start)  
    fmt.Printf("Time taken: %s\n", sortSliceTime)