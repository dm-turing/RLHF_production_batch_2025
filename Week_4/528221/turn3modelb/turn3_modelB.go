package main

import (
	"math/rand"
	"net/url"
	"time"
)

// Server represents a backend server in the pool.
type Server struct {
	URL    *url.URL
	Weight int // Weight for weighted round robin
}

// LoadBalancer interface defines the methods for a load balancer.
type LoadBalancer interface {
	Next() *url.URL
}

// SimpleLoadBalancer implements a round-robin load balancer.
type SimpleLoadBalancer struct {
	servers []*Server
	current int
}

// NewSimpleLoadBalancer creates a new round-robin load balancer.
func NewSimpleLoadBalancer(servers []*Server) LoadBalancer {
	return &SimpleLoadBalancer{servers: servers}
}

// Next implements LoadBalancer.Next().
func (lb *SimpleLoadBalancer) Next() *url.URL {
	// Round-robin logic
	lb.current = (lb.current + 1) % len(lb.servers)
	return lb.servers[lb.current].URL
}

// WeightedRoundRobinLoadBalancer implements a weighted round-robin load balancer.
type WeightedRoundRobinLoadBalancer struct {
	servers     []*Server
	totalWeight int
}

// NewWeightedRoundRobinLoadBalancer creates a new weighted round-robin load balancer.
func NewWeightedRoundRobinLoadBalancer(servers []*Server) LoadBalancer {
	// Calculate total weight for weighted round robin
	totalWeight := 0
	for _, server := range servers {
		totalWeight += server.Weight
	}
	return &WeightedRoundRobinLoadBalancer{servers: servers, totalWeight: totalWeight}
}

// Next implements LoadBalancer.Next().
func (lb *WeightedRoundRobinLoadBalancer) Next() *url.URL {
	// Weighted round-robin logic
	index := rand.Intn(lb.totalWeight)
	for _, server := range lb.servers {
		index -= server.Weight
		if index < 0 {
			return server.URL
		}
	}
	return nil // Should never happen due to totalWeight > 0
}

// Callback example for Least Connections routing
func leastConnectionsCallback(servers []*Server) *url.URL {
	// ... Least Connections logic ...
	return nil
}

func mustParseURL(ip string) *url.URL {
	url, _ := url.Parse(ip)
	return url
}

// Main function
func main() {
	rand.Seed(time.Now().UnixNano())

	// Initialize server pool
	servers := []*Server{
		{URL: mustParseURL("http://localhost:8081"), Weight: 1},
		{URL: mustParseURL("http://localhost:8082"), Weight: 2},
		{URL: mustParseURL("http://localhost:8083"), Weight: 3},
	}

	// Create load balancers using callbacks
	roundRobinLB := NewSimpleLoadBalancer(servers)
	weightedRoundRobinLB := NewWeightedRoundRobinLoadBalancer(servers)
	roundRobinLB, weightedRoundRobinLB = roundRobinLB, weightedRoundRobinLB
}
