package gloadbalancer

import "fmt"

func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
