package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type threadSafeMap struct {
	m     map[string]string
	mutex *sync.RWMutex
}

func New() *threadSafeMap {
	return &threadSafeMap{
		m:     make(map[string]string),
		mutex: &sync.RWMutex{},
	}
}

// Set stores a key-value pair in the map.
func (tm *threadSafeMap) Set(key string, value string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.m[key] = value
}

// Get retrieves a value from the map by key.
func (tm *threadSafeMap) Get(key string) (string, bool) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	value, ok := tm.m[key]
	return value, ok
}

// Delete removes a key-value pair from the map.
func (tm *threadSafeMap) Delete(key string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	delete(tm.m, key)
}

// Len returns the length of the map.
func (tm *threadSafeMap) Len() int {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return len(tm.m)
}

func main() {
	tsm := New()

	go func() {
		for i := 0; i < 1000; i++ {
			tsm.Set(strconv.Itoa(i), strconv.Itoa(i))
			value, exists := tsm.Get(strconv.Itoa(i * 10))
			if exists {
				fmt.Println("Key exists, value:", value)
			} else {
				fmt.Println("Key does not exist")
			}
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			tsm.Set(strconv.Itoa(i*10), strconv.Itoa(i*10))
			value, exists := tsm.Get(strconv.Itoa(i * 10))
			if exists {
				fmt.Println("Key exists, value:", value)
			} else {
				fmt.Println("Key does not exist")
			}
		}
	}()

	// Wait for goroutines to complete
	time.Sleep(2 * time.Second)
}
