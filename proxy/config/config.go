package config

import (
	"time"
	"runtime"
)

const (
	listenAddr    = "" // Listen to every address
	listenPort    = "1234"
	serverAddr    = "91.121.51.22" +
		"0"
	serverPort    = "27016"
	HandleTimeout = 5 * time.Second
	ServerConnectionTimeout = 1 * time.Minute

	ListenFullAddr = listenAddr + ":" + listenPort
	ServerFullAddr = serverAddr + ":" + serverPort
)

var (
	WorkerCount = runtime.NumCPU()
)
