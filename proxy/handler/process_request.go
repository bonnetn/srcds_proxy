package handler

import (
	"srcds_proxy/proxy/conntrack"
	"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/config"
	"context"
)

type requestProcessorHandler struct {
	// requestProcessorHandler forwards the client messages to the server by managing connections in the conntable.
	done            <-chan struct{}
	connectionTable conntrack.ConnectionTable
}

func NewRequestProcessorHandler(done <-chan struct{}, table conntrack.ConnectionTable) srcds.Handler {
	return &requestProcessorHandler{
		done:            done,
		connectionTable: table,
	}
}

func (h *requestProcessorHandler) Handle(
	ctx context.Context, responseWriter srcds.ConnectionWriter, msg srcds.Message, addr srcds.AddressPort) error {
	// Handle handles the incoming connections to the proxy. It forwards every byte received from the proxy to the
	// server. If the client send its first message, it creates a new connection to the server dedicated for this
	// client's traffic. This connection is added to the connection table. If this client sends messages again, the
	// handler will re-use the server connection in the connection table to forward the traffic.
	// When a new connection is created, a listener worker is also created to forward the responses back to the client.
	var (
		serverConn srcds.ConnectionWriter
		err        error
	)
	if !h.connectionTable.HasConnection(addr) {
		// If it is the first message received by this client, make a new connection to the server that will be used to
		// forward the messages from this client.
		//log.Print("New client: ", addr.String())
		udpConn, err := srcds.Dial(config.ServerFullAddr)
		if err != nil {
			return err
		}
		conn := srcds.NewConnectionWriter(*udpConn, nil)

		//log.Print("-> Create connection to ", udpConn.RemoteAddr().String())
		serverConn = h.connectionTable.GetOrStoreConnection(addr, conn)

		// Make a worker that will listen to the newly created connection and that will forward back every response.
		go func() {
			handler := NewResponseProcessorHandler(responseWriter)
			srcds.Serve(h.done, *udpConn, handler)
		}()

	} else {
		//log.Print("Known client: ", addr.String())
		serverConn, err = h.connectionTable.GetConnection(addr)
		if err != nil {
			return err
		}
	}
	return serverConn.Write(srcds.MessageToBytes(msg))
}
