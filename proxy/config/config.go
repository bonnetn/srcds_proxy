package config

import (
	"os"
	"sync"

	"github.com/golang/glog"
)

const (
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
	listenFullAddr string
	serverFullAddr string
	once           sync.Once
)

// ListenAddr returns the listen address of the proxy.
func ListenAddr() string {
	once.Do(extractConfFromEnvVars)
	return listenFullAddr
}

// ServerAddr returns the address of the SRCDS server.
func ServerAddr() string {
	once.Do(extractConfFromEnvVars)
	return serverFullAddr
}

func getEnvOrDefault(envKey string, defaultValue string) string {
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	}
	return defaultValue
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
