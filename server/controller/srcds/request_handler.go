package srcds

import (
	"net"
	"srcds_proxy/server/model/conntrack"
	"srcds_proxy/server/model/srcds_connection"
	"srcds_proxy/server/controller"
	"srcds_proxy/server/worker"
	)

type clientHandler struct {
	conntrack conntrack.ConnectionTable
	listenConn srcds_connection.SRCDSConnection
}

func NewClientHandler(conntrack conntrack.ConnectionTable, listenConn srcds_connection.SRCDSConnection) controller.Handler {
	return &clientHandler{
		conntrack: conntrack,
		listenConn: listenConn,
	}
}

func (h *clientHandler) Handle(buffer []byte, addr net.UDPAddr, n int) error {
	// Handle will handle the incoming connections to the proxy. It will forward every byte received to the server.
	// If it is a new connection, it will add an entry to the connection table and instantiate a controller that will listen
	// for responses from the server.

	clientConn, err := h.conntrack.GetConnection(addr)
	if err == conntrack.NoConnectionAssociatedError {
		if clientConn, err = h.createConnectionAndWorker(addr); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Forward all bytes.
	// TODO: add handlers.
	if _, err = clientConn.Write(buffer[0:n]); err != nil {
		return err
	}
	return nil
}

func (h *clientHandler) createConnectionAndWorker(addr net.UDPAddr) (*srcds_connection.SRCDSConnection, error) {
	// createConnectionAndWorker will create a connection in the connection table and add a listening controller.
	conn, err := h.conntrack.CreateConnection(addr)
	if err != nil {
		return nil, err
	}

	// Create a controller that will process the responses from the server.
	worker.NewProxyWorker(*conn, NewServerHandler(addr, h.listenConn))

	return conn, nil
}
