package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type EngagementMetric struct {
	ConversionRate   float64       `json:"conversion_rate"`
	PageViews        int           `json:"page_views"`
	BounceRate       float64       `json:"bounce_rate"`
	SessionDuration  time.Duration `json:"session_duration"`
	ClickThroughRate float64       `json:"click_through_rate"`
}

var (
	// Map to store metrics for each page
	pageMetrics = map[string]*EngagementMetric{}
	mux         = sync.Mutex{}
	duration    = 10 * time.Second
)

func generateMetrics() {
	for {
		mux.Lock()
		for _, metric := range pageMetrics {
			// Generate random metrics for each page
			metric.ConversionRate = rand.Float64()
			metric.PageViews = rand.Intn(1000) + 1
			metric.BounceRate = rand.Float64()
			metric.SessionDuration = time.Duration(rand.Intn(60)) * time.Second
			metric.ClickThroughRate = rand.Float64()
		}
		mux.Unlock()
		// Sleep for duration to regenerate metrics after the given time
		time.Sleep(duration)
	}
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mux.Lock()
	defer mux.Unlock()
	if page := r.FormValue("page"); page != "" {
		if metric, ok := pageMetrics[page]; ok {
			if err := json.NewEncoder(w).Encode(metric); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Page not found", http.StatusNotFound)
		}
	} else {
		if err := json.NewEncoder(w).Encode(pageMetrics); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// Initialize metrics for home and product pages
	pageMetrics["home"] = &EngagementMetric{}
	pageMetrics["product"] = &EngagementMetric{}
	// Start a goroutine to generate metrics in real time
	go generateMetrics()

	http.HandleFunc("/metrics", getMetrics)
	fmt.Println("Server is running on http://localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
