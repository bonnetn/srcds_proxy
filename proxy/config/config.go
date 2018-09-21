package config

import (
	"runtime"
	"os"
	"sync"
	"github.com/golang/glog"
	"time"
)

const (
	HandleTimeout           = 5 * time.Second
	ServerConnectionTimeout = 1 * time.Minute

	defaultListenIP   = "0.0.0.0"
	defaultListenPort = "27015"
	defaultServerIP   = "192.168.0.2"
	defaultServerPort = "27015"

	listenIPKey   = "LISTEN_ADDR"
	listenPortKey = "LISTEN_PORT"
	serverIPKey   = "SERVER_ADDR"
	serverPortKey = "SERVER_PORT"
)

var (
	workerCount    = runtime.NumCPU()
	listenFullAddr string
	serverFullAddr string
	once           sync.Once
)

const MaxDatagramSize = 4096

func ListenAddr() string {
	once.Do(extractConfFromEnvVars)
	return listenFullAddr
}

func ServerAddr() string {
	once.Do(extractConfFromEnvVars)
	return serverFullAddr
}

func WorkerCount() int {
	return workerCount
}

func getEnvOrDefault(envKey string, defaultValue string) string {
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	} else {
		return defaultValue
	}
}

func extractConfFromEnvVars() {
	var (
		listenAddr = getEnvOrDefault(listenIPKey, defaultListenIP)
		listenPort = getEnvOrDefault(listenPortKey, defaultListenPort)
		serverAddr = getEnvOrDefault(serverIPKey, defaultServerIP)
		serverPort = getEnvOrDefault(serverPortKey, defaultServerPort)
	)
	glog.Info("Extracting configuration from environment variables.")
	listenFullAddr = listenAddr + ":" + listenPort
	serverFullAddr = serverAddr + ":" + serverPort
}
