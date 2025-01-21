package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Metrics struct to hold user engagement metrics
type Metrics struct {
	Visits                int
	Clicks                int
	AddToCartEvents       int
	PurchaseCompletions   int
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
	metrics := &Metrics{}
	rand.Seed(time.Now().UnixNano())

	// Simulate 100 users going through the funnel
	for i := 0; i < 100; i++ {
		step := VisitProduct(metrics)
		for step != nil {
			step = step(metrics)
		}
	}

	// Calculate and display the metrics
	fmt.Printf("Total Visits: %d\n", metrics.Visits)
	fmt.Printf("Total Clicks: %d\n", metrics.Clicks)
	fmt.Printf("Total Add to Cart Events: %d\n", metrics.AddToCartEvents)
	fmt.Printf("Total Purchase Completions: %d\n", metrics.PurchaseCompletions)

	conversionRate := float64(metrics.PurchaseCompletions) / float64(metrics.Visits) * 100
	clickThroughRate := float64(metrics.Clicks) / float64(metrics.Visits) * 100

	fmt.Printf("Conversion Rate: %.2f%%\n", conversionRate)
	fmt.Printf("Click Through Rate: %.2f%%\n", clickThroughRate)
}
