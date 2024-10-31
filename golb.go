package golb

import (
	"log"
	"net/http"
	"sync"
	"time"
)

// LoadBalancer represents a simple round-robin load balancer with active cleaning and passive recovery.
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
	mu              sync.Mutex
}

// NewLoadBalancer initializes a LoadBalancer with the specified port and servers.
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	lb := &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
	go lb.recoverUnhealthyServers(10 * time.Second)
	return lb
}

// getNextAvailableServer returns the next available server in a round-robin fashion.
func (lb *LoadBalancer) getNextAvailableServer() Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := 0; i < len(lb.servers); i++ {
		server := lb.servers[lb.roundRobinCount%len(lb.servers)]
		if server.IsAlive() {
			lb.roundRobinCount++
			return server
		}
		lb.roundRobinCount++
	}

	log.Println("No healthy servers available")
	return nil
}

// serveProxy forwards the request to the selected server.
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	server := lb.getNextAvailableServer()
	if server == nil {
		http.Error(rw, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}
	log.Printf("Forwarding request to address %q\n", server.Address())
	server.Serve(rw, r)
}

// recoverUnhealthyServers periodically checks and recovers unhealthy servers.
func (lb *LoadBalancer) recoverUnhealthyServers(interval time.Duration) {
	for {
		time.Sleep(interval)
		lb.mu.Lock()
		for _, server := range lb.servers {
			if !server.IsAlive() {
				log.Printf("Attempting to recover server at %s", server.Address())
				server.StartHealthCheck(2 * time.Second) // Start health check on the server
			}
		}
		lb.mu.Unlock()
	}
}

// Start initializes servers, sets up the load balancer, and starts listening on the specified port.
func Start(serverAddresses []string, port string) {
	var servers []Server
	for _, addr := range serverAddresses {
		server := NewSimpleServer(addr)
		server.StartHealthCheck(2 * time.Second) // Regular health check
		servers = append(servers, server)
	}

	lb := NewLoadBalancer(port, servers)
	http.HandleFunc("/", lb.serveProxy)

	log.Printf("Load balancer listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
