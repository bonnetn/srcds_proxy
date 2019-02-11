package models

import (
	"net"
)

// UDPMaxSize is the maximum size of a UDP datagram.
const UDPMaxSize = 65507

// Packet represents a packet that was sent to the proxy.
type Packet struct {
	Src     Host
	Dst     Host
	Size    int
	Content []byte

	Socket *net.UDPConn
}
