package main

import (
	"testing"
	"time"
)

// Benchmark functions for cache operations

func BenchmarkInsert(b *testing.B) {
	cache := NewCache(100, 250, 5*time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Put(i, i)
	}
}

func BenchmarkLookup(b *testing.B) {
	cache := NewCache(100, 250, 5*time.Minute)
	for i := 0; i < 250; i++ {
		cache.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(i % 250)
	}
}

func BenchmarkDelete(b *testing.B) {
	cache := NewCache(100, 250, 5*time.Minute)
	for i := 0; i < 250; i++ {
		cache.Put(i, i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(i % 250)
	}
}

func BenchmarkMemoryAllocation(b *testing.B) {
	_ = NewCache(100, 250, 5*time.Minute)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = make([]byte, 64) // Simulate memory allocation for new cache entries
	}
}
