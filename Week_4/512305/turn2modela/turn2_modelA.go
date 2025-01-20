package main

func BubbleSort(data []int) {
	n := len(data)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

func QuickSort(data []int) {
	if len(data) < 2 {
		return
	}
	left, right := 0, len(data)-1
	pivot := data[len(data)-1]

	for i := range data {
		if data[i] < pivot {
			data[i], data[left] = data[left], data[i]
			left++
		}
	}

	data[left], data[right] = data[right], data[left]

	QuickSort(data[:left])
	QuickSort(data[left+1:])
}
