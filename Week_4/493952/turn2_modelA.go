package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define metrics
var (
	conversionRate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ecommerce_conversion_rate",
		Help: "Conversion rate",
	})
	pageViews = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ecommerce_page_views",
		Help: "Page views",
	})
	bounceRate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ecommerce_bounce_rate",
		Help: "Bounce rate",
	})
	sessionDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ecommerce_session_duration",
		Help: "Session duration",
	})
	clickThroughRate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "ecommerce_click_through_rate",
		Help: "Click-through rate",
	})
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(conversionRate, pageViews, bounceRate, sessionDuration, clickThroughRate)
}

func updateMetrics() {
	for {
		time.Sleep(10 * time.Second)
		conversionRate.Set(rand.Float64() * 100)
		pageViews.Set(float64(rand.Intn(10000)))
		bounceRate.Set(rand.Float64() * 100)
		sessionDuration.Set(float64(rand.Intn(600)))
		clickThroughRate.Set(rand.Float64() * 100)
	}
}

func main() {
	go updateMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}
