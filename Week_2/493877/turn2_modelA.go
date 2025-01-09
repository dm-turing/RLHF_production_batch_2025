package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type RateLimiter struct {
	client        *redis.Client
	windowSize    int
	maxRequests   int
	criticalRatio float64
	mu            sync.Mutex
	windowMap     map[string]int
}

func NewRateLimiter(redisURL string, windowSize int, maxRequests int, criticalRatio float64) (*RateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		client:        client,
		windowSize:    windowSize,
		maxRequests:   maxRequests,
		criticalRatio: criticalRatio,
		windowMap:     make(map[string]int),
	}, nil
}

func (rl *RateLimiter) Limit(ctx context.Context, key string, limit int, period time.Duration) bool {
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

		rl, err := NewRateLimiter("localhost:6379", 100, limit, 0.75) // 75% threshold
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

		rl, err := NewRateLimiter("localhost:6379", 100, limit, 0.75) // 75% threshold
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
