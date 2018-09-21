package main

import (
	"flag"
	"github.com/bonnetn/srcds_proxy/proxy"
)

func main() {
	flag.Parse()
	proxy.Launch()
}
