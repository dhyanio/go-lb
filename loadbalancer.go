package golb

import (
	"fmt"
	"net/http"
)

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}
func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *LoadBalancer) serverProxy(rw http.ResponseWriter, r *http.Request) {
	targetSever := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address %q\n", targetSever.Address())
	targetSever.Serve(rw, r)
}

// Start starts loadbalancer
func start(serverAdd []string, port string) {
	var servers []Server

	for _, value := range serverAdd {
		servers = append(servers, newSimpleServer(value))
	}

	lb := NewLoadBalancer(port, servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serverProxy(rw, req)
	}

	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Serving requests at localhost:%s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
