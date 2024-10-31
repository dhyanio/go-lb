package main

import (
	"github.com/dhyanio/golb"
)

func main() {
	servers := []string{"http://localhost:8081", "http://localhost:8082"}
	golb.Start(servers, "3000")
}
