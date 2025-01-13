package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	// Set the logging level to debug for detailed logs
	log.SetLevel(logrus.DebugLevel)
	// Set the formatter to output JSON logs
	log.SetFormatter(&logrus.JSONFormatter{})

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

		// Log request details
		log.WithFields(logrus.Fields{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"latency":     latency.String(),
			"remote_addr": c.Request.RemoteAddr,
			"user_agent":  c.Request.UserAgent(),
		}).Debug("Request processed")
	})

	// Define a sample handler for demonstration
	r.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Product list"})
	})

	// Start the server
	r.Run(":8080")
}
