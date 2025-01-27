package main

import (
	"sync"
	"time"
)

type Limiter struct {
	// Map of user IDs to their rate limiters
	userLimits map[string]*RateLimiter
	// Map of IP addresses to their rate limiters
	ipLimits map[string]*RateLimiter
	// Map of endpoint URLs to their rate limiters
	endpointLimits map[string]*RateLimiter
	// Synchronization mutex
	mu sync.Mutex
	// Duration of the rate limiting window
	window time.Duration
	// Maximum number of requests allowed in the window
	maxRequests int
}

// NewLimiter creates a new rate limiter with the given window duration and maximum requests.
func NewLimiter(window time.Duration, maxRequests int) *Limiter {
	return &Limiter{
		userLimits:     make(map[string]*RateLimiter),
		ipLimits:       make(map[string]*RateLimiter),
		endpointLimits: make(map[string]*RateLimiter),
		window:         window,
		maxRequests:    maxRequests,
	}
}

// RateLimiter represents a rate limiter for a specific user, IP, or endpoint.
type RateLimiter struct {
	// Time of the last request
	lastRequest time.Time
	// Number of requests made in the current window
	requestCount int
}

// Allowed checks if the request is allowed based on the rate limits.
func (l *Limiter) Allowed(userID, ip, endpoint string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	// Get or create rate limiters for the user, IP, and endpoint
	userLimiter := l.getOrCreateLimiter(l.userLimits, userID)
	ipLimiter := l.getOrCreateLimiter(l.ipLimits, ip)
	endpointLimiter := l.getOrCreateLimiter(l.endpointLimits, endpoint)
	// Check if all limiters are allowed
	return userLimiter.Allowed() && ipLimiter.Allowed() && endpointLimiter.Allowed()
}

// getOrCreateLimiter gets the rate limiter for the given key or creates a new one if it doesn't exist.
func (l *Limiter) getOrCreateLimiter(limiters map[string]*RateLimiter, key string) *RateLimiter {
	limiter, ok := limiters[key]
	if !ok {
		limiter = &RateLimiter{}
		limiters[key] = limiter
	}
	return limiter
}

// Allowed checks if the request is allowed based on the rate limit.
func (r *RateLimiter) Allowed() bool {
	now := time.Now()
	// If the window has passed, reset the counter
	if now.Sub(r.lastRequest) >= l.window {
		r.requestCount = 0
	}
	// Increment the counter
	r.requestCount++
	// Check if the limit is exceeded
	if r.requestCount > l.maxRequests {
		return false
	}
	// Update the last request time
	r.lastRequest = now
	return true
}
