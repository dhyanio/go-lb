package main

import (
	"github.com/dhyanio/golb"
)

func main() {
	serverList := []string{
		"https://www.google.com/",
		"https://github.com/dhyanio",
		"https://www.linkedin.com/in/dhyanio/",
	}

	golb.Start(serverList, "3000")
}
