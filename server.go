package golb

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// server represents a load-balanced server interface.
type server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
	StartHealthCheck(interval time.Duration)
	MarkUnhealthy()
}

// simpleServer represents a backend server with reverse proxy capabilities and health-check status.
type simpleServer struct {
	addr    string
	proxy   *httputil.ReverseProxy
	isAlive bool
	mu      sync.Mutex
}

// newSimpleServer creates a new instance of simpleServer with reverse proxy and health-check status.
func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("Failed to parse server URL %q: %v", addr, err)
	}

	return &simpleServer{
		addr:    addr,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
		isAlive: true,
	}
}

// Address returns the server address.
func (s *simpleServer) Address() string {
	return s.addr
}

// IsAlive returns the current health status of the server.
func (s *simpleServer) IsAlive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isAlive
}

// MarkUnhealthy marks the server as unhealthy.
func (s *simpleServer) MarkUnhealthy() {
	s.mu.Lock()
	s.isAlive = false
	s.mu.Unlock()
}

// Serve handles the HTTP request by forwarding it to the reverse proxy.
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// StartHealthCheck periodically checks the server’s health and updates its status.
func (s *simpleServer) StartHealthCheck(interval time.Duration) {
	go func() {
		for {
			resp, err := http.Get(s.addr + "/health")
			s.mu.Lock()
			s.isAlive = err == nil && resp.StatusCode == http.StatusOK
			s.mu.Unlock()
			time.Sleep(interval)
		}
	}()
}
