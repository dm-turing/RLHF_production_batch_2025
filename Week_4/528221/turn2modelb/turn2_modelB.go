package main

import (
	// ... (existing imports)
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"sync"
)

type ServerPool struct {
	servers         []*url.URL
	connections     map[*url.URL]int
	connectionMutex sync.Mutex
}

func NewServerPool() *ServerPool {
	p := &ServerPool{servers: make([]*url.URL, 0)}
	p.connections = make(map[*url.URL]int)
	return p
}

func (p *ServerPool) NextLeastConnections() *url.URL {
	p.connectionMutex.Lock()
	defer p.connectionMutex.Unlock()

	var leastServer *url.URL
	var leastConnections int = math.MaxInt32

	for _, server := range p.servers {
		connections := p.connections[server]
		if connections < leastConnections {
			leastServer = server
			leastConnections = connections
		}
	}
	return leastServer
}

// Increment/decrement connection count for a server
// on request start/end
func (p *ServerPool) IncrementConnection(server *url.URL) {
	p.connectionMutex.Lock()
	p.connections[server]++
	p.connectionMutex.Unlock()
}

func (p *ServerPool) DecrementConnection(server *url.URL) {
	p.connectionMutex.Lock()
	p.connections[server]--
	p.connectionMutex.Unlock()
}

func (p *ServerPool) AddServer(serverURL *url.URL) {
	p.servers = append(p.servers, serverURL)
}

func proxyHandler(p *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		server := p.NextLeastConnections()
		proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = server.Scheme
			r.URL.Host = server.Host
			r.Header.Set("X-Forwarded-Host", r.Host)
			r.Host = server.Host
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
