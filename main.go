package main

import "C2Dev24/c2"

func main() {
	listener := c2.HTTPListener{
		IP:   "127.0.0.1",
		Port: "8080",
	}

	// start listening for requests
	listener.Listen()
}
