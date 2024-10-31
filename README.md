# GoLB: Simple but easy to use Golang Loadbalancer
<img src="./doc/golb.png" alt="Discache" width="700"/>

GoLB is a simple RoundRobin golang loadbalancer. It performs active cleaning and passive recovery for unhealthy   backends.

```go
func main() {
    // Slice of servers
	servers := []string{"http://localhost:8081", "http://localhost:8082"}
    
    // Start Loadbalancer
	golb.Start(serverList, "3000")
}
```