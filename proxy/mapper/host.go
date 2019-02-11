package mapper

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/models"
)

// StringToHost converts a string like "127.0.0.1:27015" to a Host object.
func StringToHost(addr string) (*models.Host, error) {
	listenAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	listenHost, err := models.UDPAddrToHost(listenAddr)
	if err != nil {
		return nil, err
	}

	return listenHost, nil
}
