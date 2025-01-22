package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"net/url"
)

type ServerPool struct {
	servers []*url.URL
	current uint32
}

func hashIP(ip string) int {
	h := fnv.New32a()
	h.Write([]byte(ip))
	return int(h.Sum32())
}

func (p *ServerPool) NextIPHash(ip string) *url.URL {
	serverIndex := hashIP(ip) % len(p.servers)
	return p.servers[serverIndex]
}

func (p *ServerPool) AddServer(serverURL *url.URL) {
	p.NextIPHash(serverURL.String())
	p.servers = append(p.servers, serverURL)
}

func NewServerPool() *ServerPool {
	return &ServerPool{servers: make([]*url.URL, 0)}
}

func proxyHandler(p *ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverURL := p.Next()
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
