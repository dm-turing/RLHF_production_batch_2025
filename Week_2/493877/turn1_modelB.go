package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ulule/limiter"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Handle the request here
	fmt.Fprintf(w, "Welcome to the homepage!")
}

// Define a granularLimiter struct with separate limiters for different API endpoints
type granularLimiter struct {
	limiters map[string]*limiter.Limiter
}

func newGranularLimiter() *granularLimiter {
	return &granularLimiter{
		limiters: make(map[string]*limiter.Limiter),
	}
}

// getLimiter retrieves the limiter for the specified API endpoint
func (l *granularLimiter) getLimiter(endpoint string) *limiter.Limiter {
	// Determine the appropriate limiter based on the endpoint
	limiterKey := l.getLimiterKey(endpoint)

	// Check if the limiter exists
	limiter, ok := l.limiters[limiterKey]
	if !ok {
		// If the limiter does not exist, create a new one with the appropriate burst capacity
		limiter = limiter.NewRateFromConfig(limiter.RateLimitConfig{
			Rate:  1,  // Set the rate per second for each limiter
			Burst: 10, // Allow a burst of requests before rate limiting
		})
		l.limiters[limiterKey] = limiter
	}
	return limiter
}

// getLimiterKey determines the appropriate limiter key based on the endpoint path
func (l *granularLimiter) getLimiterKey(endpoint string) string {
	// You can further refine this based on your specific use case
	limiterKey := endpoint
	if strings.HasPrefix(limiterKey, "/api/v1/data/") {
		limiterKey = "/api/v1/data/*"
	}
	return limiterKey
}

func main() {
	rateLimiter := newGranularLimiter()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Apply rate limiting based on the endpoint
		limiter := rateLimiter.getLimiter(r.URL.Path)
		if limited := limiter.Allow(); !limited {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		handler(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
