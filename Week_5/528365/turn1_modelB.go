package main

import (
	"container/list"
	"fmt"
	"sync"
)

type LRUCache struct {
	capacity int
	cache    map[int]*list.Element
	list     *list.List
	mutex    sync.RWMutex
}
type entry struct {
	key   int
	value int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}
func (c *LRUCache) Get(key int) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if e, ok := c.cache[key]; ok {
		c.moveToFront(e)
		return e.Value.(*entry).value
	}
	return -1
}
func (c *LRUCache) Put(key int, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if e, ok := c.cache[key]; ok {
		e.Value.(*entry).value = value
		c.moveToFront(e)
		return
	}
	e := &entry{key: key, value: value}
	if c.list.Len() == c.capacity {
		back := c.list.Back()
		c.list.Remove(back)
		delete(c.cache, back.Value.(*entry).key)
	}
	c.list.PushFront(e)
	c.cache[key] = c.list.Front()
}
func (c *LRUCache) moveToFront(e *list.Element) {
	c.list.Remove(e)
	c.list.PushFront(e)
}
func main() {
	cache := NewLRUCache(100)
	for i := 1; i <= 250; i++ {
		cache.Put(i, i*i)
	}
	for i := 1; i <= 250; i++ {
		value := cache.Get(i)
		if value == -1 {
			fmt.Println("Key", i, "not found in the cache")
		} else {
			fmt.Println("Key", i, ":", value)
		}
	}
}
