package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ServerPool struct {
	servers []*Server
	current uint32
}

func (p *ServerPool) Next() *url.URL {
	n := atomic.AddUint32(&p.current, 1)
	return p.servers[(int(n)-1)%len(p.servers)]
}

type Server struct {
	URL    *url.URL
	Weight int // Weight for round-robin
}

func (p *ServerPool) NextWeighted() *url.URL {
	totalWeight := 0
	for _, server := range p.servers {
		totalWeight += server.Weight
	}
	index := rand.Intn(totalWeight)
	for _, server := range p.servers {
		idx -= server.Weight
		if index < 0 {
			return server.URL
		}
	}
	return nil
}

func (p *ServerPool) AddServer(serverURL *url.URL) {
	p.servers = append(p.servers, serverURL)
}

func NewServerPool() *ServerPool {
	return &ServerPool{servers: make([]*Server, 0)}
}

func proxyHandler(p *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverURL := p.NextWeighted()
		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = serverURL.Scheme
			r.URL.Host = serverURL.Host
			r.Header.Set("X-Forwarded-Host", r.Host)
			r.Host = serverURL.Host
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
