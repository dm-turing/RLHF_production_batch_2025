package main

import (
    "sync"
    "time"
)

type RateLimiter struct {
    maxRequests int
    windowTime  time.Duration
    numIntervals int
    lock        sync.Mutex
    windowStart time.Time
    requestCounts []int
}

func NewRateLimiter(maxRequests int, windowTime time.Duration) *RateLimiter {
    rl := &RateLimiter{
        maxRequests: maxRequests,
        windowTime:  windowTime,
        numIntervals: 10, // Assuming 10 intervals for simplicity
        windowStart: time.Now(),
    }
    rl.requestCounts = make([]int, rl.numIntervals)
    return rl
}

func (rl *RateLimiter) Reset() {
    rl.lock.Lock()
    defer rl.lock.Unlock()
    rl.windowStart = time.Now()
    for i := range rl.requestCounts {
        rl.requestCounts[i] = 0
    }
}

func (rl *RateLimiter) Allow() bool {
    rl.lock.Lock()
    defer rl.lock.Unlock()

    // Calculate the current window interval index
    now := time.Now()
    intervalIndex := int((now.Sub(rl.windowStart)) / (rl.windowTime / time.Duration(rl.numIntervals)))

    // Reset expired intervals
    for i := 0; i < intervalIndex; i++ {
        rl.requestCounts[i] = 0
    }

    // Increment the current interval's request count
    rl.requestCounts[intervalIndex]++

    // Check if the total request count exceeds the limit
    totalCount := 0
    for _, count := range rl.requestCounts {
        totalCount += count
    }
    return totalCount <= rl.maxRequests
}

func setMaxRequestsHandler(rl *RateLimiter) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... (Same implementation as before)
        rl.SetMaxRequests(req.MaxRequests)
        rl.Reset() // Reset the sliding window after updating max requests
        w.Write([]byte("Max requests updated successfully"))
    }
}

func setWindowTimeHandler(rl *RateLimiter) http.HandlerFunc {