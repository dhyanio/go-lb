# GoLB: Simple but easy to use Golang Loadbalancer
<img src="./doc/golb.png" alt="Discache" width="700"/>

GoLB is a simple RoundRobin golang loadbalancer. It performs active cleaning and passive recovery for unhealthy   backends.

```go
func main() {
    // Slice of servers
	serverList := []string{
		"https://www.google.com/",
		"https://github.com/dhyanio",
		"https://www.linkedin.com/in/dhyanio/",
	}

    // Start Loadbalancer
	golb.Start(serverList, "3000")
}
```