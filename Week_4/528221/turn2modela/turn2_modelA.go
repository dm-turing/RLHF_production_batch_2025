package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
)

type Server struct {
	URL         *url.URL
	ActiveConns int64 // Use atomic operations
}

type ServerPool struct {
	servers []*Server
}

func (p *ServerPool) Next() *Server {
	// Find the server with the least connections
	var selected *Server
	minConns := int64(^uint64(0) >> 1) // Max int64 value
	for _, server := range p.servers {
		activeConns := atomic.LoadInt64(&server.ActiveConns)
		if activeConns < minConns {
			minConns = activeConns
			selected = server
		}
	}
	return selected
}

func (p *ServerPool) AddServer(serverURL *url.URL) {
	selected := new(Server)
	selected.URL = serverURL
	p.servers = append(p.servers, selected)
}

func NewServerPool() *ServerPool {
	return &ServerPool{servers: make([]*Server, 0)}
}

func proxyHandler(p *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server := p.Next()
		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = server.URL.Scheme
			r.URL.Host = server.URL.Host
			r.Header.Set("X-Forwarded-Host", r.Host)
			r.Host = server.URL.Host
			http.DefaultTransport.RoundTrip(r)
		})
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	serverPool := NewServerPool()
	servers := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, server := range servers {
		url, err := url.Parse(server)
		if err != nil {
			log.Fatal(err)
		}
		serverPool.AddServer(url)
	}

	http.HandleFunc("/", proxyHandler(serverPool))
	fmt.Println("Starting load balancer on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
