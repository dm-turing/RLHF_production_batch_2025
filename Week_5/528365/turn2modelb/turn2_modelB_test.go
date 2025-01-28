package main

import (
	"testing"
)

func BenchmarkLRUCache_Insert(b *testing.B) {
	c := NewCache(100, 250, 0) // Cache without expiration
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Put(i, "test")
	}
}

func BenchmarkLRUCache_Lookup(b *testing.B) {
	c := NewCache(100, 250, 0)
	for i := 0; i < 1000; i++ {
		c.Put(i, "test")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Get(i)