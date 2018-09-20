package config

import (
	"time"
	"runtime"
	"os"
)

func getEnvOrDefault(envKey string, defaultValue string) string {
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	} else {
		return defaultValue
	}
}

var (
	listenAddr              = getEnvOrDefault("LISTEN_ADDR", "")
	listenPort              = getEnvOrDefault("LISTEN_PORT", "1234")
	serverAddr              = getEnvOrDefault("SERVER_ADDR", "91.121.51.220")
	serverPort              = getEnvOrDefault("SERVER_PORT", "27016")
	HandleTimeout           = 5 * time.Second
	ServerConnectionTimeout = 1 * time.Minute

	workerCount    = runtime.NumCPU()
	listenFullAddr = listenAddr + ":" + listenPort
	serverFullAddr = serverAddr + ":" + serverPort
)

func ListenAddr() string {
	return listenFullAddr
}

func ServerAddr() string {
	return serverFullAddr
}

func WorkerCount() int {
	return workerCount
}
