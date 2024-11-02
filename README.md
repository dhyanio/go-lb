# GoLB: The Lightweight, User-Friendly Load Balancer for Testing in Go
<img src="./doc/golb.png" alt="Discache" width="700"/>

GoLB is a simple RoundRobin golang loadbalancer. It performs active cleaning and passive recovery for unhealthy   backends.

### Features:

- `Cleans`: Excludes servers marked as unhealthy from the round-robin rotation.
- `Recovers`: Periodically tries to recheck and reintegrate previously unhealthy servers back into service.

```go
package main

import (
    "github.com/dhyanio/golb"
)

func main() {
    // Slice of servers
	servers := []string{"http://localhost:8081", "http://localhost:8082"}

    // Start Loadbalancer
	golb.Start(serverList, "3000")
}
```