package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RateLimiter struct {
	client         *redis.Client
	limits         map[string]int
	mu             sync.Mutex
	adaptionPeriod time.Duration
}

const (
	initialLimit = 100
	maxLimit     = 1000
	adaptionFreq = 5 * time.Minute // Frequency to adapt rate limits
)

func NewRateLimiter(redisURL string) (*RateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	rl := &RateLimiter{
		client:         client,
		limits:         make(map[string]int),
		mu:             sync.Mutex{},
		adaptionPeriod: adaptionFreq,
	}

	// Start a routine to adapt rate limits every adaptionFreq
	go rl.adaptRateLimits(context.Background())

	return rl, nil
}

func (rl *RateLimiter) adaptRateLimits(ctx context.Context) {
	ticker := time.NewTicker(rl.adaptionPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			for key, limit := range rl.limits {
				newLimit := rl.calculateAdaptedLimit(key, limit)
				rl.limits[key] = newLimit
			}
			rl.mu.Unlock()
		}
	}
}

// Calculates a new adapted rate limit based on current usage
func (rl *RateLimiter) calculateAdaptedLimit(key string, currentLimit int) int {
	// Simulate random traffic for this example
	traffic := rand.Intn(100)

	if traffic >= 70 { // Critical traffic threshold
		newLimit := currentLimit * 2 // Increase limit by a factor of 2 during critical traffic
		if newLimit > maxLimit {
			newLimit = maxLimit
		}
		return newLimit
	}

	// Relax rate limit if traffic is low
	newLimit := currentLimit / 2
	if newLimit < initialLimit {
		newLimit = initialLimit
	}
	return newLimit
}

func (rl *RateLimiter) Limit(ctx context.Context, key string, period time.Duration) bool {
	rl.mu.Lock()
	limit, ok := rl.limits[key]
	if !ok {
		// Initialize limit if it doesn't exist
		limit = initialLimit
		rl.limits[key] = limit
	}
	rl.mu.Unlock()
	return rl.limitRequests(ctx, key, limit, period)
}

func (rl *RateLimiter) limitRequests(ctx context.Context, key string, limit int, period time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Remove expired requests from the window
	now := time.Now()
	for k, v := range rl.windowMap {
		if time.Unix(int64(v), 0).Add(period).Before(now) {
			delete(rl.windowMap, k)
		}
	}

	// Add the current request to the window
	if _, exists := rl.windowMap[key]; !exists {
		rl.windowMap[key] = int(now.Unix())
	}

	// Count requests in the window
	requestCount := len(rl.windowMap)

	// Calculate dynamic limit based on the critical ratio
	if requestCount > rl.maxRequests*int(rl.criticalRatio) {
		// Critical traffic, be restrictive
		limit = rl.maxRequests / 2
	} else {
		// Normal traffic, be lenient
		limit = rl.maxRequests
	}

	// Check against the dynamic limit
	if requestCount >= limit {
		return true
	}

	return false
}

func main() {
	http.HandleFunc("/api/v1/resource1", func(w http.ResponseWriter, r *http.Request) {
		limit := 100
		period := time.Minute
		key := "ip:" + r.RemoteAddr

		rl, err := NewRateLimiter("localhost:6379") // 75% threshold
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rl.Limit(context.Background(), key, limit, period) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Resource1 accessed successfully\n")
	})

	http.HandleFunc("/api/v1/resource2", func(w http.ResponseWriter, r *http.Request) {
		limit := 50
		period := time.Hour
		key := "ip:" + r.RemoteAddr

		rl, err := NewRateLimiter("localhost:6379") // 75% threshold
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if rl.Limit(context.Background(), key, limit, period) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Resource2 accessed successfully\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
