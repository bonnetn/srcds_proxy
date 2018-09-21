package main

import (
	"flag"
	"srcds_proxy/proxy"
	)

func main() {
	flag.Parse()
	proxy.Launch()
}
