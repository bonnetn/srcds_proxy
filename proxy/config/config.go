package config

import (
	"time"
	"runtime"
)

const (
	listenAddr     = "" // Listen to every address
	listenPort     = "1234"
	serverAddr     = "127.0.0.1"
	serverPort     = "27016"
	HandleTimeout = 5 * time.Second

	ListenFullAddr = listenAddr + ":" + listenPort
	ServerFullAddr = serverAddr + ":" + serverPort
)

var (
	WorkerCount = runtime.NumCPU()
)