package srcds

import (
	"net"
	"fmt"
)

type AddressPort struct {
	// AddressPort is the IP address+port of a client. This can be used as a key in a map.
	port uint16
	ip   [4]uint8
}

func UDPAddrToAddressPort(udpAddr net.UDPAddr) AddressPort {
	ip := udpAddr.IP.To4()
	return AddressPort{
		port: uint16(udpAddr.Port),
		ip:   [4]uint8{ip[0], ip[1], ip[2], ip[3]},
	}
}

func (a *AddressPort) getIP() net.IP {
	return net.IPv4(a.ip[0], a.ip[1], a.ip[2], a.ip[3])
}

func (a *AddressPort) getPort() int {
	return int(a.port)
}

func (a *AddressPort) String() string {
	return fmt.Sprintf("%s:%d", a.getIP().String(), a.getPort())

}
