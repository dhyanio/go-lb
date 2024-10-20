package main

import (
	"fmt"
	"net/http"
)

func main() {
	servers := []Server{
		newSimpleServer("https://www.kubefront.com"),
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.linkedin.com/in/dhyanio/"),
	}
	lb := NewLoadBalancer("8080", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serverProxy(rw, req)

	}
	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Serving requests at localhost:%s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
