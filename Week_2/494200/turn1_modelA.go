package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger middleware using logrus
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture request information
		logrus.Infof("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()

		// Capture response information
		statusCode := c.Writer.Status()
		logrus.Infof("Response sent: status %d, latency %s", statusCode, gin.StatusDuration(c))

		// Capture potential security alerts like authentication failures
		if c.Errors.Has() {
			for _, err := range c.Errors {
				logrus.Warnf("Error: %s", err.Error())
			}
		}
	}
}

// An example handler for a GET request
func getProducts(c *gin.Context) {
	products := []string{"Laptop", "Smartphone", "Tablet"}
	c.JSON(http.StatusOK, products)
}

func main() {
	// Initialize logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Create a new gin router
	router := gin.Default()

	// Add logging middleware
	router.Use(loggerMiddleware())

	// Add routes
	router.GET("/products", getProducts)

	// Start server
	router.Run(":8080")
}
