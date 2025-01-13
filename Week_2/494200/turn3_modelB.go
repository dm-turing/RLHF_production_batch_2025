package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()

	// Metrics for logging
	logCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status_code"},
	)

	logLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func main() {
	// Register metrics with Prometheus
	prometheus.MustRegister(logCount, logLatency)

	// Set up Logrus to export metrics to Prometheus
	logrusPrometheus := prometheus.New(logCount, logLatency)
	log.AddHook(logrusPrometheus)

	// Set the logging level to debug for detailed logs
	log.SetLevel(logrus.DebugLevel)
	// Set the formatter to output JSON logs
	log.SetFormatter(&logrus.JSONFormatter{})

	// Initialize the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a new Gin router
	r := gin.Default()

	// Define a logging middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()

		// Process the request
		c.Next()

		// Record the end time
		end := time.Now()
		latency := end.Sub(start)

		// Check if we need to log this request
		if shouldLog(c) {
			// Log request details
			log.WithFields(logrus.Fields{
				"method":      c.Request.Method,
				"path":        c.Request.URL.Path,
				"status":      c.Writer.Status(),
				"latency":     latency.String(),
				"remote_addr": c.Request.RemoteAddr,
				"user_agent":  c.Request.UserAgent(),
			}).Debug("Request processed")
		}
	})

	// Define a sample handler for demonstration
	r.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Product list"})
	})

	// Define a high-traffic handler with log sampling
	r.GET("/product/:id", func(c *gin.Context) {
		productID := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Product with ID %s", productID)})
	})

	// Start the server
	r.Run(":8080")
}

// shouldLog checks if a request should be logged based on a sampling rate
func shouldLog(c *gin.Context) bool {
	// Set sampling rate (log 1 out of every 100 requests)
	samplingRate := 100

	// Generate a random number
	randomInt := rand.Intn(samplingRate)

	// Log this request if the random number is zero
	return randomInt == 0
}
