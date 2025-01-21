package main

import (
	"image/color"
	"math/rand"
	"runtime"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type BenchmarkResult struct {
	size          int
	executionTime float64
	memUsage      uint64
	comparisons   int
	swaps         int
}

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

var mcomparisons *int

func MergeSort(arr []int) (int, int) {
	mcomparisons = new(int)
	*mcomparisons = 0
	MergeSortHelper(arr)
	return *mcomparisons, 0
}

// Merge sort implementation
func MergeSortHelper(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	left := MergeSortHelper(arr[:mid])
	right := MergeSortHelper(arr[mid:])

	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0)
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		*mcomparisons++
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

func measureBenchmark(sortFunc func([]int) (int, int), size int) BenchmarkResult {
	data := generateRandomSlice(size)
	memStatsStart := new(runtime.MemStats)
	runtime.ReadMemStats(memStatsStart)

	startTime := time.Now()
	comparisons, swaps := sortFunc(data)
	elapsed := time.Since(startTime).Seconds()

	memStatsEnd := new(runtime.MemStats)
	runtime.ReadMemStats(memStatsEnd)

	memUsage := memStatsEnd.TotalAlloc - memStatsStart.TotalAlloc

	return BenchmarkResult{
		size:          size,
		executionTime: elapsed,
		memUsage:      memUsage,
		comparisons:   comparisons,
		swaps:         swaps,
	}
}

func createBarChart(results []BenchmarkResult, attribute string) {
	p := plot.New()
	p.Title.Text = "Execution Time / Memory Usage"
	p.Y.Label.Text = "Value"
	w := vg.Points(10)

	bars := make(plotter.Values, len(results))
	cats := make([]string, len(results))

	for i, result := range results {
		var value float64
		switch attribute {
		case "time":
			value = result.executionTime
			p.Y.Label.Text = "Execution Time (Seconds)"
		case "memory":
			value = float64(result.memUsage) / 1024
			p.Y.Label.Text = "Memory Usage (KB)"
		}
		bars[i] = value
		cats[i] = string(result.size)
	}

	b, err := plotter.NewBarChart(bars, w)
	b.Color = color.RGBA{R: 2, A: 255}
	if err != nil {
		panic(err)
	}

	p.Add(b)
	p.NominalX(cats...)

	err = p.Save(4*vg.Inch, 4*vg.Inch, attribute+"_bar_chart.png")
	if err != nil {
		panic(err)
	}
}

func createLineGraph(results []BenchmarkResult) {
	p := plot.New()

	p.Title.Text = "Execution Time vs Input Size"
	p.X.Label.Text = "Size of Input"
	p.Y.Label.Text = "Execution Time (Seconds)"

	pts := make(plotter.XYs, len(results))

	for i, result := range results {
		pts[i].X = float64(result.size)
		pts[i].Y = result.executionTime
	}

	l, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}

	p.Add(l)

	err = p.Save(4*vg.Inch, 4*vg.Inch, "time_line_graph.png")
	if err != nil {
		panic(err)
	}
}

func createScatterPlot(results []BenchmarkResult) {
	p := plot.New()
	p.Title.Text = "Comparisons vs Swaps"
	p.X.Label.Text = "Comparisons"
	p.Y.Label.Text = "Swaps"

	pts := make(plotter.XYs, len(results))

	for i, result := range results {
		pts[i].X = float64(result.comparisons)
		pts[i].Y = float64(result.swaps)
	}

	s, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}

	p.Add(s)

	err = p.Save(4*vg.Inch, 4*vg.Inch, "comparison_swaps_scatter_plot.png")
	if err != nil {
		panic(err)
	}
}

func main() {
	inputSizes := []int{10, 100, 1000}
	results := []BenchmarkResult{}

	// Measure results
	for _, size := range inputSizes {
		results = append(results, measureBenchmark(BubbleSort, size))
		results = append(results, measureBenchmark(MergeSort, size))
		results = append(results, measureBenchmark(HeapSort, size))
	}

	// Create charts for different metrics
	createBarChart(results, "time")
	createBarChart(results, "memory")
	createLineGraph(results)
	createScatterPlot(results)
}
