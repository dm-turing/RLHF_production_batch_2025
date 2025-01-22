package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	URL     string
	IsAlive bool
}

type LoadBalancer struct {
	Servers []Server
	mu      sync.Mutex
	index   int
}

func (lb *LoadBalancer) NextServer() *Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.Servers) == 0 {
		return nil
	}

	lb.index = (lb.index + 1) % len(lb.Servers)
	return &lb.Servers[lb.index]
}
func (lb *LoadBalancer) IsAlive() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for _, server := range lb.Servers {
		if server.IsAlive {
			return true
		}
	}
	return false
}

func (lb *LoadBalancer) healthCheck() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		lb.mu.Lock()
		for i, server := range lb.Servers {
			_, err := http.Get(server.URL)
			if err != nil {
				lb.Servers[i].IsAlive = false
			} else {
				lb.Servers[i].IsAlive = true
			}
		}
		lb.mu.Unlock()
	}
}

func main() {
	lb := LoadBalancer{
		Servers: []Server{
			{URL: "http://server1:8000", IsAlive: true},
			{URL: "http://server2:8000", IsAlive: true},
		},
	}

	go lb.healthCheck()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := lb.NextServer()
		if server == nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, "No available servers")
			return
		}

		// Forward the request to the selected server
		resp, err := http.Get(server.URL + r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprintln(w, err)
			return
		}

		defer resp.Body.Close()

		// Copy the response from the server to the client
		w.WriteHeader(resp.StatusCode)
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			fmt.Println("Error copying response:", err)
		}
	})

	fmt.Println("Load balancer running on port 8080")
	http.ListenAndServe(":8080", nil)
}
