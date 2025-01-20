package turn1modelb

func BubbleSort(data []int) []int {
	n := len(data)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
	return data
}

// Merge function to merge two sorted slices
func Merge(left, right []int) []int {
	result := []int{}
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

	// Append remaining elements
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// MergeSort function to sort an array
func MergeSort(data []int) []int {
	if len(data) < 2 {
		return data // Base case
	}

	mid := len(data) / 2
	left := MergeSort(data[:mid])
	right := MergeSort(data[mid:])

	return Merge(left, right)
}

// InsertionSort function to sort an array
func InsertionSort(data []int) []int {
	for i := 1; i < len(data); i++ {
		key := data[i]
		j := i - 1

		// Move elements that are greater than key to one position ahead
		for j >= 0 && data[j] > key {
			data[j+1] = data[j]
			j--
		}
		data[j+1] = key
	}
	return data
}

// SelectionSort function to sort an array
func SelectionSort(data []int) []int {
	for i := 0; i < len(data)-1; i++ {
		minIndex := i
		for j := i + 1; j < len(data); j++ {
			if data[j] < data[minIndex] {
				minIndex = j
			}
		}
		// Swap the found minimum element with the first element
		data[i], data[minIndex] = data[minIndex], data[i]
	}
	return data
}

// QuickSort function to sort an array
func QuickSort(data []int) []int {
	if len(data) < 2 {
		return data // Base case
	}

	pivot := data[len(data)/2] // Choosing the pivot element
	left := []int{}
	right := []int{}
	equal := []int{}

	for _, value := range data {
		if value < pivot {
			left = append(left, value)
		} else if value > pivot {
			right = append(right, value)
		} else {
			equal = append(equal, value)
		}
	}

	return append(append(QuickSort(left), equal...), QuickSort(right)...)
}

// Heapify function to maintain the heap property
func Heapify(data []int, n, i int) {
	largest := i     // Initialize largest as root
	left := 2*i + 1  // left = 2*i + 1
	right := 2*i + 2 // right = 2*i + 2

	// If left child is larger than root
	if left < n && data[left] > data[largest] {
		largest = left
	}

	// If right child is larger than largest so far
	if right < n && data[right] > data[largest] {
		largest = right
	}

	// If largest is not root
	if largest != i {
		data[i], data[largest] = data[largest], data[i] // Swap
		Heapify(data, n, largest)                       // Recursively heapify the affected subtree
	}
}

// HeapSort function to sort an array
func HeapSort(data []int) []int {
	n := len(data)

	// Build heap (rearrange array)
	for i := n/2 - 1; i >= 0; i-- {
		Heapify(data, n, i)
	}

	// One by one extract elements from heap
	for i := n - 1; i > 0; i-- {
		data[0], data[i] = data[i], data[0] // Swap
		Heapify(data, i, 0)
	}

	return data
}

// BucketSort function to sort an array
func BucketSort(data []int) []int {
	bucketSize := 5
	if len(data) == 0 {
		return data
	}

	// Find the maximum value in data
	maxVal := data[0]
	for _, value := range data {
		if value > maxVal {
			maxVal = value
		}
	}

	// Create buckets
	bucketCount := (maxVal / bucketSize) + 1
	buckets := make([][]int, bucketCount)

	// Distribute the elements into buckets
	for _, value := range data {
		index := value / bucketSize
		buckets[index] = append(buckets[index], value)
	}

	// Sort individual buckets and concatenate
	output := []int{}
	for _, bucket := range buckets {
		if len(bucket) > 0 {
			// Use any sorting algorithm; using Insertion Sort here
			output = append(output, InsertionSort(bucket)...)
		}
	}

	return output
}

// GetMax function to find the maximum number in the array
func GetMax(data []int) int {
	maxVal := data[0]
	for _, value := range data {
		if value > maxVal {
			maxVal = value
		}
	}
	return maxVal
}
func CountingSort(data []int) []int {
	if len(data) == 0 {
		return data
	}

	// Step 1: Find the range of the input data
	maxVal := data[0]
	for _, num := range data {
		if num > maxVal {
			maxVal = num
		}
	}

	// Step 2: Create a count array
	count := make([]int, maxVal+1)

	// Step 3: Store counts of each number
	for _, num := range data {
		count[num]++
	}

	// Step 4: Reconstruct the sorted array
	sortedIndex := 0
	for i, c := range count {
		for c > 0 {
			data[sortedIndex] = i
			sortedIndex++
			c--
		}
	}

	return data
}

func getMax(data []int) int {
	maxVal := data[0]
	for _, num := range data {
		if num > maxVal {
			maxVal = num
		}
	}
	return maxVal
}

func countingSortForRadix(data []int, exp int) {
	n := len(data)
	output := make([]int, n)
	count := make([]int, 10)

	// Count occurrences of each digit
	for i := 0; i < n; i++ {
		index := (data[i] / exp) % 10
		count[index]++
	}

	// Change count[i] so that it contains the actual position of
	// this digit in output[]
	for i := 1; i < 10; i++ {
		count[i] += count[i-1]
	}

	// Build the output array
	for i := n - 1; i >= 0; i-- {
		index := (data[i] / exp) % 10
		output[count[index]-1] = data[i]
		count[index]--
	}

	// Copy the output array to data, so that data now
	// contains sorted numbers according to the current digit
	for i := 0; i < n; i++ {
		data[i] = output[i]
	}
}

func RadixSort(data []int) []int {
	maxVal := getMax(data)

	// Apply counting sort to sort elements based on significant digits
	for exp := 1; maxVal/exp > 0; exp *= 10 {
		countingSortForRadix(data, exp)
	}

	return data
}
