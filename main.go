package main

import (
	"C2Dev24/c2"
)

func main() {
	listener := c2.HTTPListener{
		IP:   "127.0.0.1",
		Port: "8080",
	}
	// Run listener as go routine
	go listener.Listen()

	// Start CLI
	c2.StartCLI()
}
