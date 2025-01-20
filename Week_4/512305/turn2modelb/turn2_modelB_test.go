package main

import (
	"testing"
)

type SorterTestFixture struct {
	OriginalData []int
	TestData     []int
	parent       *SorterTestFixture
}

// NewSorterTestFixture creates a new SorterTestFixture with a parent
func NewSorterTestFixture(parent *SorterTestFixture) *SorterTestFixture {
	return &SorterTestFixture{parent: parent}
}

// Setup initializes resources for the test
func (f *SorterTestFixture) Setup(data []int) {
	f.OriginalData = data
	f.TestData = make([]int, len(data))
	copy(f.TestData, data)

	if f.parent != nil {
		f.parent.Setup(data)
	}
}

// Teardown releases resources or resets them
func (f *SorterTestFixture) Teardown() {
	f.TestData = nil

	if f.parent != nil {
		f.parent.Teardown()
	}
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

// NewNestedFixture creates a nested SorterTestFixture
func (f *SorterTestFixture) NewNestedFixture() *SorterTestFixture {
	return NewSorterTestFixture(f)
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
	data := []int{64, 34, 25, 12, 22, 11, 90}
	expected := []int{11, 12, 22, 25, 34, 64, 90}

	// Initialize top-level fixture
	fixture := NewSorterTestFixture(&SorterTestFixture{
		OriginalData: data,
		TestData:     expected,
		parent:       nil,
	})

	// Setup for a nested fixture
	nestedFixture := fixture.NewNestedFixture()

	// Test BubbleSort with nested fixture
	nestedFixture.RunTest(t, "BubbleSort", BubbleSort, expected)

	// Test QuickSort with the top-level fixture (nested setup/teardown will also be called)
	fixture.RunTest(t, "QuickSort", QuickSort, expected)
}
