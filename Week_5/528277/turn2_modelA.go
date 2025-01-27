package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests      int
	windowTime       time.Duration
	lock             sync.Mutex
	userRequests     map[string]int
	ipRequests       map[string]int
	endpointRequests map[string]int
	lastAccess       map[string]time.Time
}

func NewRateLimiter(maxRequests int, windowTime time.Duration) *RateLimiter {
	return &RateLimiter{
		maxRequests:      maxRequests,
		windowTime:       windowTime,
		userRequests:     make(map[string]int),
		ipRequests:       make(map[string]int),
		endpointRequests: make(map[string]int),
		lastAccess:       make(map[string]time.Time),
	}
}

func (rl *RateLimiter) Allow(user, ip, endpoint string) bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	now := time.Now()
	rl.cleanup(now)

	userKey := fmt.Sprintf("user:%s", user)
	ipKey := fmt.Sprintf("ip:%s", ip)
	endpointKey := fmt.Sprintf("endpoint:%s", endpoint)

	if rl.userRequests[userKey] >= rl.maxRequests ||
		rl.ipRequests[ipKey] >= rl.maxRequests ||
		rl.endpointRequests[endpointKey] >= rl.maxRequests {
		return false
	}

	rl.userRequests[userKey]++
	rl.ipRequests[ipKey]++
	rl.endpointRequests[endpointKey]++
	rl.lastAccess[userKey] = now
	rl.lastAccess[ipKey] = now
	rl.lastAccess[endpointKey] = now

	return true
}

func (rl *RateLimiter) cleanup(now time.Time) {
	for key, lastTime := range rl.lastAccess {
		if now.Sub(lastTime) > rl.windowTime {
			if key[:5] == "user:" {
				delete(rl.userRequests, key)
			}
			if key[:3] == "ip:" {
				delete(rl.ipRequests, key)
			}
			if key[:9] == "endpoint:" {
				delete(rl.endpointRequests, key)
			}
			delete(rl.lastAccess, key)
		}
	}
}

func (rl *RateLimiter) SetMaxRequests(maxRequests int) {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	rl.maxRequests = maxRequests
}

func (rl *RateLimiter) SetWindowTime(windowTime time.Duration) {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	rl.windowTime = windowTime
}

func handler(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Header.Get("X-User-ID")
		ip := r.RemoteAddr
		endpoint := r.URL.Path

		if !rl.Allow(user, ip, endpoint) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.Write([]byte("Request successful!"))
	}
}

func setMaxRequestsHandler(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type Request struct {
			MaxRequests int `json:"max_requests"`
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rl.SetMaxRequests(req.MaxRequests)
		w.Write([]byte("Max requests updated successfully"))
	}
}

func setWindowTimeHandler(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type Request struct {
			WindowTime int `json:"window_time"` // Expected in seconds for simplicity
		}
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rl.SetWindowTime(time.Duration(req.WindowTime) * time.Second)
		w.Write([]byte("Window time updated successfully"))
	}
}

func main() {
	rateLimiter := NewRateLimiter(100, time.Minute)

	http.HandleFunc("/", handler(rateLimiter))
	http.HandleFunc("/setmaxrequests", setMaxRequestsHandler(rateLimiter))
	http.HandleFunc("/setwindowtime", setWindowTimeHandler(rateLimiter))

	http.ListenAndServe(":8080", nil)
}
