package golb

import (
	"log"
	"net/http"
	"sync"
)

// LoadBalancer represents a simple round-robin load balancer.
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
	mu              sync.Mutex
}

// NewLoadBalancer initializes a LoadBalancer with the specified port and servers.
func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

// getNextAvailableServer returns the next available server in a round-robin fashion.
func (lb *LoadBalancer) getNextAvailableServer() Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

// serveProxy forwards the request to the selected server.
func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	log.Printf("Forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(rw, r)
}

// Start initializes servers, sets up the load balancer, and starts listening on the specified port.
func Start(serverAddresses []string, port string) {
	var servers []Server

	for _, addr := range serverAddresses {
		servers = append(servers, NewSimpleServer(addr))
	}

	lb := NewLoadBalancer(port, servers)
	http.HandleFunc("/", lb.serveProxy)

	log.Printf("Load balancer listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
