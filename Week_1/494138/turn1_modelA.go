package main

import (
	"fmt"
	"sync"
	"time"
)

type ThreadSafeMap struct {
	mu sync.RWMutex
	m  map[string]int
}

func NewThreadSafeMap() *ThreadSafeMap {
	return &ThreadSafeMap{
		m: make(map[string]int),
	}
}

func (tsm *ThreadSafeMap) Set(key string, value int) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	tsm.m[key] = value
}

func (tsm *ThreadSafeMap) Get(key string) (int, bool) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	return tsm.m[key], tsm.m[key] != 0
}

func (tsm *ThreadSafeMap) Delete(key string) {
	tsm.mu.Lock()
	defer tsm.mu.Unlock()
	delete(tsm.m, key)
}

func main() {
	tsm := NewThreadSafeMap()

	go func() {
		for i := 0; i < 1000; i++ {
			tsm.Set("key", i)
			value, exists := tsm.Get("key")
			if exists {
				fmt.Println("Key exists, value:", value)
			} else {
				fmt.Println("Key does not exist")
			}
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			tsm.Set("key", i*10)
			value, exists := tsm.Get("key")
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
