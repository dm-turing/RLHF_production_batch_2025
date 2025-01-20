package main

import (
	"testing"
)

// SorterTestFixture for managing test fixtures
type SorterTestFixture struct {
	OriginalData []int
	TestData     []int
}

// Setup initializes resources for the test
func (f *SorterTestFixture) Setup(data []int) {
	f.OriginalData = data
	f.TestData = make([]int, len(data))
	copy(f.TestData, data)
}

// Teardown releases resources or resets them
func (f *SorterTestFixture) Teardown() {
	f.TestData = nil
}

// RunTest executes the test with setup and teardown
func (f *SorterTestFixture) RunTest(t *testing.T, name string, sort func([]int), expected []int) {
	f.Setup(f.OriginalData)
	defer f.Teardown()

	sort(f.TestData)

	if !equal(f.TestData, expected) {
		t.Errorf("%s failed: expected %v, got %v", name, expected, f.TestData)
	}
}

// Helper function to check equality of two slices
func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// SortingSuiteFixture manages a suite of sorting tests
type SortingSuiteFixture struct {
	DataSets    [][]int
	Expected    [][]int
	SubFixtures []*SorterTestFixture
}

// Setup initializes resources for the entire suite
func (ssf *SortingSuiteFixture) Setup() {
	ssf.SubFixtures = make([]*SorterTestFixture, len(ssf.DataSets))
	for i, data := range ssf.DataSets {
		subFixture := &SorterTestFixture{}
		subFixture.Setup(data)
		ssf.SubFixtures[i] = subFixture
	}
}

// Teardown releases resources for the entire suite
func (ssf *SortingSuiteFixture) Teardown() {
	for _, subFixture := range ssf.SubFixtures {
		subFixture.Teardown()
	}
	ssf.SubFixtures = nil
}

// RunSuiteTest executes the suite of tests
func (ssf *SortingSuiteFixture) RunSuiteTest(t *testing.T, sortFuncs []func([]int)) {
	ssf.Setup()
	defer ssf.Teardown()

	for i, subFixture := range ssf.SubFixtures {
		for _, sortFunc := range sortFuncs {
			subFixture.RunTest(t, funcName(sortFunc), sortFunc, ssf.Expected[i])
		}
	}
}

// Helper function to get function name for labeling
func funcName(i interface{}) string {
	switch i.(type) {
	case func([]int):
		return "Sorting Function"
	default:
		return "Unknown"
	}
}

func TestNestedSortingSuite(t *testing.T) {
	// Define different datasets and the expected order
	dataSets := [][]int{
		{5, 2, 9, 1, 5, 6},
		{3, 0, -1, 8, 7},
	}
	expected := [][]int{
		{1, 2, 5, 5, 6, 9},
		{-1, 0, 3, 7, 8},
	}

	// Initialize the suite fixture
	suiteFixture := &SortingSuiteFixture{
		DataSets: dataSets,
		Expected: expected,
	}

	// Define sorting functions to test
	sortFuncs := []func([]int){BubbleSort, QuickSort}

	// Run the suite
	suiteFixture.RunSuiteTest(t, sortFuncs)
}
