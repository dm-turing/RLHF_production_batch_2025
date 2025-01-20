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

func TestSortingAlgorithms(t *testing.T) {
	// Test case data
	data := []int{5, 2, 9, 1, 5, 6}
	expected := []int{1, 2, 5, 5, 6, 9}

	// Initialize test fixture
	fixture := &SorterTestFixture{}

	// Setup
	fixture.Setup(data)

	// Test BubbleSort
	fixture.RunTest(t, "BubbleSort", BubbleSort, expected)

	// Test QuickSort
	fixture.RunTest(t, "QuickSort", QuickSort, expected)

	// Teardown
	fixture.Teardown()
}
