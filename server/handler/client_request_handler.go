package handler

import (
	"net"
	"srcds_proxy/server/conntrack"
)

type clientHandler struct {
	conntrack conntrack.ConntrackTable
}

func NewClientHandler(conntrack conntrack.ConntrackTable) Handler {
	return &clientHandler{
		conntrack: conntrack,
	}
}

func (h *clientHandler) Handle(buffer []byte, addr *net.UDPAddr, n int) error {
	// Get or create a connection to the server.
	clientCon, err := h.conntrack.GetOrCreateConnection(*addr)
	if err != nil {
		return err
	}

	// Forward all bytes.
	// TODO: add handlers.
	_, err = clientCon.Write(buffer[0:n])
	if err != nil {
		return err
	}
}
