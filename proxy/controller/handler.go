package controller

import "net"

type Handler interface {
	Handle(buffer []byte, addr net.UDPAddr, n int) error
}
