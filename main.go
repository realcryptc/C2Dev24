package main

import (
	"C2Dev24/c2"
	"time"
)

func main() {
	listener := c2.HTTPListener{
		IP:   "127.0.0.1",
		Port: "8080",
	}

	// Create Agent
	c2.AgentMap.Add(&c2.Agent{ID: "Yo Momma", IP: "127.0.0.1", LastCall: time.Now()})
	c2.AgentMap.Add(&c2.Agent{ID: "Will Im", IP: "127.0.0.1", LastCall: time.Now()})
	c2.AgentMap.Add(&c2.Agent{ID: "IU Admin", IP: "127.0.0.2", LastCall: time.Now()})

	// start listening for requests
	go listener.Listen()

	// start our CLI
	c2.StartCLI()
}
