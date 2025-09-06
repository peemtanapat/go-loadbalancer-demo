package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	URL     *url.URL `json:"url"`
	Healthy bool     `json:"healthy"`
	mutex   sync.RWMutex
}

type LoadBalancer struct {
	servers []Server
	current uint64
}

type HealthCheckResponse struct {
	Status    string    `json:"status"`
	Instance  string    `json:"instance"`
	Port      string    `json:"port"`
	Timestamp time.Time `json:"timestamp"`
}

type StatusResponse struct {
	LoadBalancer string    `json:"loadBalancer"`
	Servers      []Server  `json:"servers"`
	Algorithm    string    `json:"algorithm"`
	Timestamp    time.Time `json:"timestamp"`
}

func (s *Server) SetHealth(healthy bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Healthy = healthy
}

func (s *Server) IsHealthy() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Healthy
}

func NewLoadBalancer() *LoadBalancer {
	servers := []Server{
		{URL: parseURL("http://host.docker.internal:8081"), Healthy: true},
		{URL: parseURL("http://host.docker.internal:8082"), Healthy: true},
		{URL: parseURL("http://host.docker.internal:8083"), Healthy: true},
	}

	return &LoadBalancer{
		servers: servers,
		current: 0,
	}
}

// Round-robin algorithm
func (lb *LoadBalancer) GetNextServer() (*Server, error) {
	healthyServers := []*Server{}

	for i := range lb.servers {
		if lb.servers[i].IsHealthy() {
			healthyServers = append(healthyServers, &lb.servers[i])
		}
	}

	if len(healthyServers) == 0 {
		return nil, fmt.Errorf("no healthy servers available")
	}

	next := atomic.AddUint64(&lb.current, 1)
	return healthyServers[next%uint64(len(healthyServers))], nil
}

func (lb *LoadBalancer) HealthCheck() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for {
		log.Println("Performing health checks (/health) to each server")

		for i := range lb.servers {
			server := &lb.servers[i]

			res, err := client.Get(server.URL.String() + "/health")
			wasHealthy := server.IsHealthy()

			if err != nil {
				server.SetHealth(false)
				if wasHealthy {
					log.Printf("❌ Server %s health check failed: %v", server.URL.String(), err)
				}
				continue
			}

			res.Body.Close()

			healthy := res.StatusCode == http.StatusOK
			server.SetHealth(healthy)

			if !wasHealthy && healthy {
				log.Printf("✅ Server %s is back up", server.URL.String())
			} else if wasHealthy && !healthy {
				log.Printf("❌ Server %s is down", server.URL.String())
			} else {
				log.Printf("...Server %s is still up", server.URL.String())
			}
		}

		time.Sleep(30 * time.Second)
	}
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/lb-status" {
		lb.handleStatus(w, r)
		return
	}

	server, err := lb.GetNextServer()
	if err != nil {
		http.Error(w, "Service Unavailable: " + err.Error(), http.StatusServiceUnavailable)
		return
	}

	log.Printf("Routing request to %s", server.URL.String())

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(server.URL)

	// Custom error handler
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
				log.Printf("❌ Proxy error for %s: %v", server.URL.String(), err)
		server.SetHealth(false)
		http.Error(w, "Service Temporarily Unavailable", http.StatusServiceUnavailable)
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		log.Printf("✅ Request completed: %s -> %d", server.URL.String(), resp.StatusCode)
		return nil
	}

	proxy.ServeHTTP(w, r)
}

func (lb *LoadBalancer) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := StatusResponse{
		LoadBalancer: "active",
		Servers: lb.servers,
		Algorithm: "round-robin",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func main() {
	lb := NewLoadBalancer()

	// Health checking in background
	go lb.HealthCheck()

	router := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("📥 [%s] %s %s", time.Now().Format("15:04:05"), r.Method, r.URL.Path)

		lb.ServeHTTP(w, r)

		log.Printf("⏱️  Request completed in %v", time.Since(startTime))
	})

	port := "9080"

	fmt.Printf("🚀 Go Load Balancer starting on port %s\n", port)
	fmt.Printf("🔍 Status endpoint: http://localhost:%s/lb-status\n", port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func parseURL(rawURL string) *url.URL {
	url, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}
	return url
}
