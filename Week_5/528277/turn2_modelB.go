package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	limits     map[string]*limit
	windowTime time.Duration
	lock       sync.Mutex
}

type limit struct {
	maxRequests int
	requests    int
	lastAccess  time.Time
}

func NewRateLimiter(windowTime time.Duration) *RateLimiter {
	return &RateLimiter{
		limits:     make(map[string]*limit),
		windowTime: windowTime,
	}
}

func (rl *RateLimiter) Allow(endpoint string) bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	l, ok := rl.limits[endpoint]
	if !ok {
		l = &limit{}
		rl.limits[endpoint] = l
	}
	now := time.Now()
	rl.cleanup(now)
	if l.requests >= l.maxRequests {
		return false
	}
	l.requests++
	l.lastAccess = now

	return true
}

func (rl *RateLimiter) SetLimit(endpoint string, maxRequests int) {
	rl.lock.Lock()
	defer rl.lock.Unlock()
	rl.limits[endpoint] = &limit{
		maxRequests: maxRequests,
	}
}

func (rl *RateLimiter) cleanup(now time.Time) {
	for _, l := range rl.limits {
		if now.Sub(l.lastAccess) > rl.windowTime {
			l.requests = 0
		}
	}
}

type handler struct {
	rl *RateLimiter
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path

	// Check if the rate limit for the endpoint is set
	if _, ok := h.rl.limits[endpoint]; !ok {
		http.Error(w, "Rate limit not configured for this endpoint", http.StatusForbidden)
		return
	}

	if !h.rl.Allow(endpoint) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	w.Write([]byte("Request successful!"))
}

func setLimitHandler(rl *RateLimiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		endpoint := r.URL.Query().Get("endpoint")
		maxRequestsStr := r.URL.Query().Get("maxRequests")
		maxRequests, err := strconv.Atoi(maxRequestsStr)
		if err != nil {
			http.Error(w, "Invalid maxRequests parameter", http.StatusBadRequest)
			return
		}

		rl.SetLimit(endpoint, maxRequests)
		fmt.Fprintf(w, "Rate limit for %s set to %d requests\n", endpoint, maxRequests)
	}
}

func main() {
	rateLimiter := NewRateLimiter(time.Minute)

	h := &handler{rl: rateLimiter}
	http.Handle("/", h)

	http.HandleFunc("/set-limit", setLimitHandler(rateLimiter))

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
