package main

import (
	"srcds_proxy/proxy"
	"log"
)

func main() {
	for {
		proxy.Launch()
		log.Print("Proxy crashed, restarting...")
	}
}
