package srcds

import (
	"net"
	"srcds_proxy/server/controller"
		"srcds_proxy/server/model/srcds_connection"
	)

type serverHandler struct {
	clientAddr net.UDPAddr
	listenConn srcds_connection.SRCDSConnection
}

func NewServerHandler(clientAddr net.UDPAddr, listenConn srcds_connection.SRCDSConnection) controller.Handler {
	return &serverHandler{
		clientAddr: clientAddr,
		listenConn: listenConn,
	}
}

func (h *serverHandler) Handle(buffer []byte, addr net.UDPAddr, n int) error {
	if _, err := h.listenConn.WriteToUDP(buffer[:n], &h.clientAddr); err != nil {
		return err
	}
	return nil
}
