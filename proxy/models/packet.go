package models

import (
	"net"

	"github.com/pkg/errors"
)

// UDPMaxSize is the maximum size of a UDP datagram.
const UDPMaxSize = 65507

// HostToUDPAddr converts a Host object into a net.UDPAddr object.
func HostToUDPAddr(host *Host) *net.UDPAddr {
	port := int(host.Port)
	addr := net.UDPAddr{
		IP:   net.IPv4(host.IP[0], host.IP[1], host.IP[2], host.IP[3]),
		Port: port,
	}
	return &addr

}

// UDPAddrToHost converts a UDPAddr object into a Host object.
func UDPAddrToHost(addr *net.UDPAddr) (*Host, error) {
	ip := addr.IP.To4()
	if ip == nil {
		return nil, errors.New("addr must be an IPv4")
	}
	if addr.Port >= 65536 { // Port can be stored with 16 bits.
		return nil, errors.New("port must be < 65536")
	}
	addrBytes := [4]byte{ip[0], ip[1], ip[2], ip[3]}
	port := uint16(addr.Port)
	host := Host{addrBytes, port}
	return &host, nil

}

// Packet represents a packet that was sent to the proxy.
type Packet struct {
	Src     Host
	Dst     Host
	Size    int
	Content []byte

	Socket *net.UDPConn
}
