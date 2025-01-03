package main

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	key   string
	value string
	next  *entry
}

type bucket struct {
	header *entry
	lock   sync.RWMutex
}

type threadSafeMap struct {
	buckets   []*bucket
	size      int
	threshold int
	mutex     sync.RWMutex
}

const (
	defaultLoadFactor float64 = 0.75
	minBuckets                = 8
)

func New() *threadSafeMap {
	tsm := &threadSafeMap{
		buckets:   make([]*bucket, minBuckets),
		size:      0,
		threshold: int(float64(minBuckets) * defaultLoadFactor),
		mutex:     sync.RWMutex{},
	}

	for i := range tsm.buckets {
		tsm.buckets[i] = &bucket{header: nil}
	}

	return tsm
}

func (tm *threadSafeMap) hash(key string) int {
	return hashString(key) % len(tm.buckets)
}

func hashString(s string) int {
	h := 0
	for _, c := range s {
		h = 31*h + int(c)
	}
	return h
}

func (tm *threadSafeMap) rehash() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	oldBuckets := tm.buckets
	tm.buckets = make([]*bucket, len(oldBuckets)*2)

	for i := range oldBuckets {
		oldBuckets[i].lock.RLock()
		e := oldBuckets[i].header
		for ; e != nil; e = e.next {
			newIndex := tm.hash(e.key)
			tm.buckets[newIndex].lock.Lock()
			newEntry := &entry{key: e.key, value: e.value}
			newEntry.next = tm.buckets[newIndex].header
			tm.buckets[newIndex].header = newEntry
			tm.buckets[newIndex].lock.Unlock()
			tm.size++
		}
		oldBuckets[i].lock.RUnlock()
	}

	tm.threshold = int(float64(len(tm.buckets)) * defaultLoadFactor)
}

// Set stores a key-value pair in the map.
func (tm *threadSafeMap) Set(key, value string) {
	tm.mutex.RLock()
	if tm.size >= tm.threshold {
		tm.rehash()
	}
	tm.mutex.RUnlock()

	index := tm.hash(key)
	tm.buckets[index].lock.Lock()
	e := tm.buckets[index].header

	for ; e != nil; e = e.next {
		if e.key == key {
			e.value = value
			tm.buckets[index].lock.Unlock()
			return
		}
	}

	newEntry := &entry{key: key, value: value}
	newEntry.next = tm.buckets[index].header
	tm.buckets[index].header = newEntry
	tm.size++
	tm.buckets[index].lock.Unlock()
}

// Get retrieves a value from the map by key.
func (tm *threadSafeMap) Get(key string) (string, bool) {
	index := tm.hash(key)
	tm.buckets[index].lock.RLock()
	e := tm.buckets[index].header

	for ; e != nil; e = e.next {
		if e.key == key {
			tm.buckets[index].lock.RUnlock()
			return e.value, true
		}
	}

	tm.buckets[index].lock.RUnlock()
	return "", false
}

// Delete removes a key-value pair from the map.
func (tm *threadSafeMap) Delete(key string) {
	index := tm.hash(key)
	tm.buckets[index].lock.Lock()
	e := tm.buckets[index].header

	var prev *entry
	for ; e != nil; e = e.next {
		if e.key == key {
			if prev == nil {
				tm.buckets[index].header = e.next
			} else {
				prev.next = e.next
			}
			tm.size--
			tm.buckets[index].lock.Unlock()
			return
		}
		prev = e
	}

	tm.buckets[index].lock.Unlock()
}

// Len returns the length of the map.
func (tm *threadSafeMap) Len() int {
	tm.mutex.RLock()
	size := tm.size
	tm.mutex.RUnlock()
	return size
}

func main() {
	tsm := New()

	go func() {
		for i := 0; i < 10; i++ {
			tsm.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			value, exists := tsm.Get(fmt.Sprintf("key%d", i))
			if exists {
				fmt.Printf("Get: %s -> %s\n", fmt.Sprintf("key%d", i), value)
			} else {
				fmt.Printf("Key %s not found\n", fmt.Sprintf("key%d", i))
			}
		}
	}()

	time.Sleep(1 * time.Second)
}
