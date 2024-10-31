package golb

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Server represents a load-balanced server interface.
type Server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
	StartHealthCheck(interval time.Duration)
	MarkUnhealthy()
}

// SimpleServer represents a backend server with reverse proxy capabilities and health-check status.
type SimpleServer struct {
	addr    string
	proxy   *httputil.ReverseProxy
	isAlive bool
	mu      sync.Mutex
}

// NewSimpleServer creates a new instance of SimpleServer with reverse proxy and health-check status.
func NewSimpleServer(addr string) *SimpleServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatalf("Failed to parse server URL %q: %v", addr, err)
	}

	return &SimpleServer{
		addr:    addr,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
		isAlive: true,
	}
}

// Address returns the server address.
func (s *SimpleServer) Address() string {
	return s.addr
}

// IsAlive returns the current health status of the server.
func (s *SimpleServer) IsAlive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isAlive
}

// MarkUnhealthy marks the server as unhealthy.
func (s *SimpleServer) MarkUnhealthy() {
	s.mu.Lock()
	s.isAlive = false
	s.mu.Unlock()
}

// Serve handles the HTTP request by forwarding it to the reverse proxy.
func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

// StartHealthCheck periodically checks the server’s health and updates its status.
func (s *SimpleServer) StartHealthCheck(interval time.Duration) {
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
