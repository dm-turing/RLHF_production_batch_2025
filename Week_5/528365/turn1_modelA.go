package main

import (
	"container/list"
	"fmt"
)

type Cache struct {
	capacity  int
	maxLength int
	cache     map[interface{}]*list.Element
	records   *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

// NewCache initializes an LRU cache with a specific capacity.
func NewCache(capacity, maxLength int) *Cache {
	return &Cache{
		capacity:  capacity,
		maxLength: maxLength,
		cache:     make(map[interface{}]*list.Element),
		records:   list.New(),
	}
}

// Get retrieves the value of the given key.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	if element, found := c.cache[key]; found {
		c.records.MoveToFront(element)
		return element.Value.(*entry).value, true
	}
	return nil, false
}

// Put adds a key-value pair to the cache.
func (c *Cache) Put(key interface{}, value interface{}) {
	if element, found := c.cache[key]; found {
		c.records.MoveToFront(element)
		element.Value.(*entry).value = value
	} else {
		entry := &entry{key, value}
		element := c.records.PushFront(entry)
		c.cache[key] = element

		if c.records.Len() > c.maxLength {
			c.removeOldest()
		} else if c.records.Len() > c.capacity && c.records.Len() <= c.maxLength {
			c.removeOldest()
		}
	}
}

// RemoveOldest removes the oldest entry from the cache.
func (c *Cache) removeOldest() {
	if element := c.records.Back(); element != nil {
		c.records.Remove(element)
		entry := element.Value.(*entry)
		delete(c.cache, entry.key)
	}
}

func main() {
	cache := NewCache(100, 250)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Put(3, "three")

	if value, found := cache.Get(2); found {
		fmt.Println("Cache hit!", value) // Should print "Cache hit! two"
	} else {
		fmt.Println("Cache miss!")
	}

	cache.Put(4, "four")
	cache.Put(5, "five")

	// Going beyond the designated capacity of 100
	for i := 6; i <= 260; i++ {
		cache.Put(i, fmt.Sprintf("%d", i))
	}

	// The first few entries will be removed as per LRU policy.
	if value, found := cache.Get(1); found {
		fmt.Println("Cache hit!", value)
	} else {
		fmt.Println("Cache miss!") // Should print "Cache miss!"
	}
}
