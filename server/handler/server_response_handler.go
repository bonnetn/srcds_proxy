package handler

import (
	"net"
	"srcds_proxy/server/conntrack"
)

type serverHandler struct {
	clientAddr net.UDPAddr
}

func NewServerHandler(clientAddr net.UDPAddr) Handler {
	return &serverHandler{
		clientAddr: clientAddr,
	}
}

func (h *serverHandler) Handle(buffer []byte, addr *net.UDPAddr, n int) error {
}
