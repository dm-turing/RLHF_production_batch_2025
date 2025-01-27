package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type SlidingWindowRateLimiter struct {
	maxRequests int
	windowTime  time.Duration
	lock        sync.Mutex
	requestLogs map[string][]time.Time
}

func NewSlidingWindowRateLimiter(maxRequests int, windowTime time.Duration) *SlidingWindowRateLimiter {
	return &SlidingWindowRateLimiter{
		maxRequests: maxRequests,
		windowTime:  windowTime,
		requestLogs: make(map[string][]time.Time),
	}
}

func (rl *SlidingWindowRateLimiter) Allow(key string) bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	now := time.Now()
	logs := rl.requestLogs[key]

	// Remove timestamps that are outside of the window
	var toKeep []time.Time
	for _, t := range logs {
		if now.Sub(t) < rl.windowTime {
			toKeep = append(toKeep, t)
		}
	}
	rl.requestLogs[key] = toKeep

	if len(toKeep) < rl.maxRequests {
		rl.requestLogs[key] = append(rl.requestLogs[key], now)
		return true
	}

	return false
}

func handler(rl *SlidingWindowRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr // Simplified; consider more attributes to form a unique key.

		if !rl.Allow(key) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		w.Write([]byte("Request successful!"))
	}
}

// Handlers to modify maxRequests and windowTime
func setMaxRequestsHandler(rl *SlidingWindowRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			MaxRequests int `json:"max_requests"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rl.lock.Lock()
		rl.maxRequests = req.MaxRequests
		rl.lock.Unlock()
		w.Write([]byte("Max requests updated successfully"))
	}
}

func setWindowTimeHandler(rl *SlidingWindowRateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			WindowTime int `json:"window_time"` // Expected in seconds
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rl.lock.Lock()
		rl.windowTime = time.Duration(req.WindowTime) * time.Second
		rl.lock.Unlock()
		w.Write([]byte("Window time updated successfully"))
	}
}

func main() {
	rateLimiter := NewSlidingWindowRateLimiter(100, time.Minute)

	http.HandleFunc("/", handler(rateLimiter))
	http.HandleFunc("/setmaxrequests", setMaxRequestsHandler(rateLimiter))
	http.HandleFunc("/setwindowtime", setWindowTimeHandler(rateLimiter))

	http.ListenAndServe(":8080", nil)
}
