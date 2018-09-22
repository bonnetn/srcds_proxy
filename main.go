package main

import (
	"flag"

	"github.com/bonnetn/srcds_proxy/proxy"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	if err := proxy.Launch(); err != nil {
		glog.Error("Failed to launch: ", err)
	}
}
