package mapper

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/pkg/errors"
)

// StringToHost converts a string like "127.0.0.1:27015" to a Host object.
func StringToHost(addr string) (*models.Host, error) {
	listenAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	listenHost, err := UDPAddrToHost(listenAddr)
	if err != nil {
		return nil, err
	}

	return listenHost, nil
}

// HostToUDPAddr converts a Host object into a net.UDPAddr object.
func HostToUDPAddr(host *models.Host) *net.UDPAddr {
	port := int(host.Port)
	addr := net.UDPAddr{
		IP:   net.IPv4(host.IP[0], host.IP[1], host.IP[2], host.IP[3]),
		Port: port,
	}
	return &addr

}

// UDPAddrToHost converts a UDPAddr object into a Host object.
func UDPAddrToHost(addr *net.UDPAddr) (*models.Host, error) {
	ip := addr.IP.To4()
	if ip == nil {
		return nil, errors.New("addr must be an IPv4")
	}
	if addr.Port >= 65536 { // Port can be stored with 16 bits.
		return nil, errors.New("port must be < 65536")
	}
	addrBytes := [4]byte{ip[0], ip[1], ip[2], ip[3]}
	port := uint16(addr.Port)
	host := models.Host{addrBytes, port}
	return &host, nil

}
