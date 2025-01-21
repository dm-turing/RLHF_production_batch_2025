package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// UserMetrics stores individual user metrics
type UserMetrics struct {
	UserID              int
	Visits              int
	Clicks              int
	AddToCartEvents     int
	PurchaseCompletions int
}

// AnalyticsChannel is used to simulate sending data to an external analytics tool
var AnalyticsChannel = make(chan UserMetrics, 100)
var wg sync.WaitGroup

// FunnelStep represents a function that processes part of the funnel
type FunnelStep func(*UserMetrics) FunnelStep

// VisitProduct simulates a user visiting a product page
func VisitProduct(metrics *UserMetrics) FunnelStep {
	return func(nextMetrics *UserMetrics) FunnelStep {
		metrics.Visits++
		fmt.Printf("User %d visited the product page.\n", metrics.UserID)
		return ClickProduct(metrics)
	}
}

// ClickProduct simulates a user clicking on a product
func ClickProduct(metrics *UserMetrics) FunnelStep {
	return func(nextMetrics *UserMetrics) FunnelStep {
		if rand.Intn(100) < 70 {
			metrics.Clicks++
			fmt.Printf("User %d clicked on the product.\n", metrics.UserID)
			return AddToCart(metrics)
		}
		return nil
	}
}

// AddToCart simulates a user adding a product to the cart
func AddToCart(metrics *UserMetrics) FunnelStep {
	return func(nextMetrics *UserMetrics) FunnelStep {
		if rand.Intn(100) < 50 {
			metrics.AddToCartEvents++
			fmt.Printf("User %d added the product to the cart.\n", metrics.UserID)
			return CompletePurchase(metrics)
		}
		return nil
	}
}

// CompletePurchase simulates a user completing a purchase
func CompletePurchase(metrics *UserMetrics) FunnelStep {
	return func(nextMetrics *UserMetrics) FunnelStep {
		if rand.Intn(100) < 30 {
			metrics.PurchaseCompletions++
			fmt.Printf("User %d completed the purchase.\n", metrics.UserID)
		}
		// Send metrics to the AnalyticsChannel
		AnalyticsChannel <- *metrics
		return nil
	}
}

// AnalyticsCollector simulates processing and displaying collected metrics
func AnalyticsCollector(wg *sync.WaitGroup) {
	for metrics := range AnalyticsChannel {
		fmt.Printf("\n--- Metrics for User %d ---\n", metrics.UserID)
		fmt.Printf("Visits: %d\n", metrics.Visits)
		fmt.Printf("Clicks: %d\n", metrics.Clicks)
		fmt.Printf("Add to Cart Events: %d\n", metrics.AddToCartEvents)
		fmt.Printf("Purchase Completions: %d\n", metrics.PurchaseCompletions)
		fmt.Println("-------------------------\n")
		wg.Done()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const numUsers = 5 // Simulate metrics for 5 users

	// Launch the analytics collector as a Goroutine
	wg.Add(numUsers)
	go AnalyticsCollector(&wg)

	// Simulate each user going through the funnel
	for userID := 1; userID <= numUsers; userID++ {
		metrics := &UserMetrics{UserID: userID}
		step := VisitProduct(metrics)
		for step != nil {
			step = step(metrics)
		}
	}

	// Allow time for all metrics to be processed
	wg.Wait()
	close(AnalyticsChannel)
}
