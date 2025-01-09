package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ContentFilterService defines the interface for the content filter service.
type ContentFilterService interface {
	FilterContent(content string) (string, error)
}

// DefaultContentFilterService is the default implementation of the ContentFilterService.
type DefaultContentFilterService struct {
	BlacklistedWords []string
}

// FilterContent filters the given content by removing blacklisted words.
func (s *DefaultContentFilterService) FilterContent(content string) (string, error) {
	filteredContent := content
	for _, word := range s.BlacklistedWords {
		filteredContent = strings.ReplaceAll(filteredContent, word, "")
	}
	return filteredContent, nil
}

// NewContentFilterService returns a new instance of the ContentFilterService.
func NewContentFilterService(blacklistedWords []string) ContentFilterService {
	return &DefaultContentFilterService{BlacklistedWords: blacklistedWords}
}

// CallbackResponse represents the structure of the dynamic API gateway response.
type CallbackResponse struct {
	StatusCode int         `json:"statusCode"`
	Body       string      `json:"body"`
	Headers    http.Header `json:"headers"`
}

// ContentFilterHandler is the HTTP handler for the content filter service.
func ContentFilterHandler(c *gin.Context, service ContentFilterService, callback func(CallbackResponse)) {
	// Extract the content from the request body.
	content := c.PostForm("content")

	// Filter the content using the service.
	filteredContent, err := service.FilterContent(content)
	if err != nil {
		log.Printf("Error filtering content: %v", err)
		callback(CallbackResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error",
			Headers:    http.Header{},
		})
		return
	}

	// Create a response with the filtered content.
	callback(CallbackResponse{
		StatusCode: http.StatusOK,
		Body:       filteredContent,
		Headers:    http.Header{},
	})
}

func main() {
	// Create a new Gin router.
	router := gin.Default()

	// Create a new instance of the content filter service with blacklisted words.
	blacklistedWords := []string{"bad", "word"}
	contentFilterService := NewContentFilterService(blacklistedWords)

	// Define the API endpoint for content filtering.
	router.POST("/filter", func(c *gin.Context) {
		// Define a callback function to handle the dynamic API gateway response.
		callback := func(response CallbackResponse) {
			c.Status(response.StatusCode)
			c.JSON(response.StatusCode, gin.H{
				"body":    response.Body,
				"headers": response.Headers,
			})
		}

		// Call the content filter handler with the callback function.
		ContentFilterHandler(c, contentFilterService, callback)
	})

	// Start the Gin router on port 8080.
	router.Run(":8080")
}
