package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/segmentio/analytics-go"
)

// Define your Segment.io write key
const segmentWriteKey = "YOUR_SEGMENT_WRITE_KEY"

// Metrics struct to hold user engagement metrics
type Metrics struct {
	UserID              string
	Visits              int
	Clicks              int
	AddToCartEvents     int
	PurchaseCompletions int
}

// FunnelStep represents a function that processes part of the funnel
type FunnelStep func(*Metrics) FunnelStep

// VisitProduct simulates a user visiting a product page
func VisitProduct(metrics *Metrics) FunnelStep {
	return func(nextMetrics *Metrics) FunnelStep {
		metrics.Visits++
		fmt.Println("User visited the product page.")
		return ClickProduct(metrics)
	}
}

// ClickProduct simulates a user clicking on a product
func ClickProduct(metrics *Metrics) FunnelStep {
	return func(nextMetrics *Metrics) FunnelStep {
		// Simulate a 70% chance of clicking
		if rand.Intn(100) < 70 {
			metrics.Clicks++
			fmt.Println("User clicked on the product.")
			return AddToCart(metrics)
		}
		return nil
	}
}

// AddToCart simulates a user adding a product to the cart
func AddToCart(metrics *Metrics) FunnelStep {
	return func(nextMetrics *Metrics) FunnelStep {
		// Simulate a 50% chance of adding to cart
		if rand.Intn(100) < 50 {
			metrics.AddToCartEvents++
			fmt.Println("User added the product to the cart.")
			return CompletePurchase(metrics)
		}
		return nil
	}
}

// CompletePurchase simulates a user completing a purchase
func CompletePurchase(metrics *Metrics) FunnelStep {
	return func(nextMetrics *Metrics) FunnelStep {
		// Simulate a 30% chance of purchase completion
		if rand.Intn(100) < 30 {
			metrics.PurchaseCompletions++
			fmt.Println("User completed the purchase.")
		}
		return nil
	}
}

func main() {
	// Initialize Segment client
	client := analytics.New(segmentWriteKey)

	// Simulate user engagement for a few users
	users := []string{"user1", "user2", "user3"}
	for _, user := range users {
		metrics := &Metrics{UserID: user}
		rand.Seed(time.Now().UnixNano())

		// Simulate 50-100 user visits for each user
		numVisits := rand.Intn(51) + 50
		for i := 0; i < numVisits; i++ {
			step := VisitProduct(metrics)
			for step != nil {
				step = step(metrics)
			}
		}

		// Send analytics events for each user
		client.Enqueue(analytics.Track{Type: "User Engagement", Properties: analytics.Properties{
			"UserID":             metrics.UserID,
			"Total Visits":       metrics.Visits,
			"Total Clicks":       metrics.Clicks,
			"Total Add to Cart":  metrics.AddToCartEvents,
			"Total Purchases":    metrics.PurchaseCompletions,
			"Conversion Rate":    float64(metrics.PurchaseCompletions) / float64(metrics.Visits) * 100,
			"Click Through Rate": float64(metrics.Clicks) / float64(metrics.Visits) * 100,
		}})
	}

	fmt.Println("Analytics events sent to Segment.")
}
