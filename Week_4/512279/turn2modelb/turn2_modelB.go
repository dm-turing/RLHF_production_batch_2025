package main

import (
	"testing"
	"time"
)

// ... (existing imports)

// ... (existing code)

func main() {
	b.Run("BubbleSort", func(b *testing.B) {
		// ... (existing bubble sort benchmark code)
		bubbleSortTimes[size] = time.Since(start)
		bubbleSortComparisons[size] = comparisons
		bubbleSortSwaps[size] = swaps
	})
	// ... (Similarly, run benchmarks for MergeSort and HeapSort)

	visualizeMetrics()
}

func visualizeMetrics() {
	// Create plots for different metrics
	createBarChart("Execution Time (ns)", "Data Size", []string{"BubbleSort", "MergeSort", "HeapSort"}, []time.Duration{bubbleSortTimes[10], mergeSortTimes[10], heapSortTimes[10]}, []time.Duration{bubbleSortTimes[1000], mergeSortTimes[1000], heapSortTimes[1000]}, []time.Duration{bubbleSortTimes[10000], mergeSortTimes[10000], heapSortTimes[10000]})

	createLineGraph("Time Complexity vs Data Size", "Data Size", "Time (ns)", []string{"BubbleSort", "MergeSort", "HeapSort"}, []time.Duration{bubbleSortTimes[10], bubbleSortTimes[1000], bubbleSortTimes[10000]}, []time.Duration{mergeSortTimes[10], mergeSortTimes[1000], mergeSortTimes[10000]}, []time.Duration{heapSortTimes[10], heapSortTimes[1000], heapSortTimes[10000]})

	createScatterPlot("Comparisons vs Swaps", "Comparisons", "Swaps", bubbleSortComparisons, bubbleSortSwaps, mergeSortComparisons, mergeSortSwaps, heapSortComparisons, heapSortSwaps)
}

// ... (helper functions to create each type of plot)

func createBarChart(title string, xLabel string, labels []string, bar1Values []time.Duration, bar2Values []time.Duration, bar3Values []time.Duration) {
	// ... (Plot creation logic)
}

func createLineGraph(title string, xLabel string, yLabel string, labels []string, line1Values []time.Duration, line2Values []time.Duration, line3Values []time.Duration) {
	// ... (Plot creation logic)
}

func createScatterPlot(title string, xLabel string, yLabel string, xValues1, yValues1, xValues2, yValues2, xValues3, yValues3 map[int]int) {
	// ... (Plot creation logic)
}

func createBarChart(title string, xLabel string, labels []string, bar1Values []time.Duration, bar2Values []time.Duration, bar3Values []time.Duration) {
    p, err := plot.New()
    if err != nil {
        panic(err)
    }
    p.Title.Text = title
    p.X.Label.Text = xLabel
    w := vg.Points(50)

    bars1 := plotter.Values(make([]float64, len(bar1Values)))