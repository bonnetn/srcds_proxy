package main

import (
	"srcds_proxy/server"
	"log"
)

func main() {
	for {
		server.Launch()
		log.Print("Proxy crashed, restarting...")
	}
}
