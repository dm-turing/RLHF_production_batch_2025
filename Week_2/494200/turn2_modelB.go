package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()
var sampleRate = float32(0.01) // 1% sample rate

func main() {
	// Set the logging level to debug for detailed logs
	log.SetLevel(logrus.DebugLevel)
	// Set the formatter to output JSON logs
	log.SetFormatter(&logrus.JSONFormatter{})

	// Create a new Gin router
	r := gin.Default()

	// Define a logging middleware with sampling
	r.Use(func(c *gin.Context) {
		start := time.Now()
		// Function to check if a request should be sampled
		shouldSample := func(endpoint string) bool {
			if endpoint == "/products/:id" {
				return rand.Float32() < sampleRate
			}
			return true
		}

		// Process the request
		c.Next()

		// Record the end time and calculate latency
		end := time.Now()
		latency := end.Sub(start)

		// extract endpoint from path
		endpoint := c.Request.URL.Path
		if idx := strings.LastIndex(endpoint, "/"); idx > 0 {
			endpoint = endpoint[:idx+1]
		}

		if shouldSample(endpoint) {
			// Log request details for sampled requests
			log.WithFields(logrus.Fields{
				"method":      c.Request.Method,
				"path":        endpoint,
				"status":      c.Writer.Status(),
				"latency":     latency.String(),
				"remote_addr": c.Request.RemoteAddr,
				"user_agent":  c.Request.UserAgent(),
			}).Debug("Request processed")
		}
	})

	// Define a sample handler for demonstration
	r.GET("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Product ID: %s", productID)})
	})

	// Define another sample handler for demonstration with less traffic
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Start the server
	r.Run(":8080")
}
