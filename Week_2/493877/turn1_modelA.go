package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(redisURL string) (*RateLimiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "",
		DB:       0,
	})

	success, err := client.Ping().Result()
	fmt.Println(success)
	if err != nil {
		return nil, err
	}

	return &RateLimiter{client: client}, nil
}

func (rl *RateLimiter) Limit(ctx context.Context, key string, limit int, period time.Duration) bool {
	val, err := rl.client.Get(key).Result()
	if val == "" {
		success, _ := rl.client.Set(key, fmt.Sprintf("%d", 1), period).Result()
		if success != "OK" {
			return true
		}
	}
	if err != nil {
		return true
	} else {
		val, _ = rl.client.Get(key).Result()
	}
	count, err := strconv.Atoi(val)
	if err != nil {
		return true
	}
	if count >= limit {
		return true
	}
	return false
}

func main() {
	http.HandleFunc("/api/v1/resource1", func(w http.ResponseWriter, r *http.Request) {
		limit := 100
		period := time.Minute
		key := "ip:" + r.RemoteAddr

		rl, err := NewRateLimiter("localhost:6379")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !rl.Limit(context.Background(), key, limit, period) {
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

		rl, err := NewRateLimiter("localhost:6379")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !rl.Limit(context.Background(), key, limit, period) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Resource2 accessed successfully\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
