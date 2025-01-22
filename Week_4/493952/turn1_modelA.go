package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// Metrics struct to hold the random metrics
type Metrics struct {
	ConversionRate   float64 `json:"conversion_rate"`
	PageViews        int     `json:"page_views"`
	BounceRate       float64 `json:"bounce_rate"`
	SessionDuration  int     `json:"session_duration"`
	ClickThroughRate float64 `json:"click_through_rate"`
}

var (
	currentMetrics Metrics
	mu             sync.RWMutex
)

func generateRandomMetrics() Metrics {
	return Metrics{
		ConversionRate:   rand.Float64() * 100,
		PageViews:        rand.Intn(10000),
		BounceRate:       rand.Float64() * 100,
		SessionDuration:  rand.Intn(600), // in seconds
		ClickThroughRate: rand.Float64() * 100,
	}
}

func updateMetrics() {
	for {
		time.Sleep(10 * time.Second)
		mu.Lock()
		currentMetrics = generateRandomMetrics()
		mu.Unlock()
	}
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentMetrics)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	currentMetrics = generateRandomMetrics()

	go updateMetrics()

	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":8080", nil)
}
