package turn1modelb

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

// TestFixture represents a test fixture for sorting algorithms.
type TestFixture struct {
	t        *testing.T
	unsorted []int
	sorted   []int
}

// NewTestFixture creates a new TestFixture with a random unsorted array.
func NewTestFixture(t *testing.T, size int) *TestFixture {
	tf := &TestFixture{t: t}
	tf.unsorted = make([]int, size)
	tf.sorted = make([]int, size)
	for i := 0; i < size; i++ {
		tf.unsorted[i] = rand.Intn(100)
	}
	copy(tf.sorted, tf.unsorted)
	sort.Ints(tf.sorted)
	return tf
}

// RunTest runs the specified sorting algorithm and asserts that the result is correct.
func (tf *TestFixture) RunTest(name string, sortAlgorithm func([]int) []int) {
	tf.t.Run(name, func(t *testing.T) {
		copy(tf.unsorted, tf.sorted)
		sortAlgorithm(tf.unsorted)
		if !reflect.DeepEqual(tf.unsorted, tf.sorted) {
			t.Errorf("Expected %v, got %v", tf.sorted, tf.unsorted)
		}
	})
}

func TestBubbleSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("BubbleSort", BubbleSort)
}

func TestInsertionSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("InsertionSort", InsertionSort)
}

func TestSelectionSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("SelectionSort", SelectionSort)
}

func TestMergeSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("MergeSort", MergeSort)
}

func TestQuickSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("QuickSort", QuickSort)
}

func TestHeapSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("HeapSort", HeapSort)
}

func TestCountingSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("CountingSort", CountingSort)
}

func TestBucketSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("BucketSort", BucketSort)
}

func TestRadixSort(t *testing.T) {
	tf := NewTestFixture(t, 10)
	tf.RunTest("RadixSort", RadixSort)
}
