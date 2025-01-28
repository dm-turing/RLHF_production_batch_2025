package main

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	mu        sync.RWMutex
	capacity  int
	maxLength int
	cache     map[interface{}]*list.Element
	records   *list.List
	expiry    time.Duration
}

type entry struct {
	key    interface{}
	value  interface{}
	expire time.Time
}

// NewCache initializes an LRU cache with a specific capacity, maxLength, and optional expiry duration.
// An expiry of 0 means no expiration.
func NewCache(capacity, maxLength int, expiry time.Duration) *Cache {
	return &Cache{
		capacity:  capacity,
		maxLength: maxLength,
		cache:     make(map[interface{}]*list.Element),
		records:   list.New(),
		expiry:    expiry,
	}
}

// Get retrieves the value of the given key.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	element, found := c.cache[key]
	if !found {
		return nil, false
	}

	entry := element.Value.(*entry)
	// Check if the entry is expired
	if entry.expire.Before(time.Now()) {
		c.mu.RUnlock()
		c.Delete(key)
		c.mu.RLock()
		return nil, false
	}

	c.records.MoveToFront(element)
	return entry.value, true
}

// Put adds a key-value pair to the cache with optional expiration.
func (c *Cache) Put(key interface{}, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Remove the old entry if exists
	if element, found := c.cache[key]; found {
		c.records.Remove(element)
	}

	expire := time.Now().Add(c.expiry)
	entry := &entry{key, value, expire}
	element := c.records.PushFront(entry)
	c.cache[key] = element

	// Evict the oldest entry if needed
	if c.records.Len() > c.maxLength {
		c.removeOldest()
	} else if c.records.Len() > c.capacity {
		c.removeOldest()
	}
}

// Delete removes a key-value pair from the cache.
func (c *Cache) Delete(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		c.records.Remove(element)
		delete(c.cache, key)
	}
}

// removeOldest removes the oldest entry from the cache.
func (c *Cache) removeOldest() {
	element := c.records.Back()
	if element != nil {
		c.records.Remove(element)
		entry := element.Value.(*entry)
		delete(c.cache, entry.key)
	}
}
