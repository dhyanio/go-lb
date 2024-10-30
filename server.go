package golb

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Server represents a load-balanced server interface.
type Server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
}

// SimpleServer represents a basic HTTP server with reverse proxy capabilities.
type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

// NewSimpleServer creates a new instance of SimpleServer with a reverse proxy.
func NewSimpleServer(addr string) *SimpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)

	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

// Address returns the server address.
func (s *SimpleServer) Address() string {
	return s.addr
}

// IsAlive always returns true, but this could be expanded for real health checks.
func (s *SimpleServer) IsAlive() bool {
	return true
}

// Serve handles the HTTP request by forwarding it to the reverse proxy.
func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}
