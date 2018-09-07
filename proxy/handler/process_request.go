package handler

import (
	"srcds_proxy/proxy/conntrack"
	"srcds_proxy/proxy/srcds"
			"srcds_proxy/proxy/config"
)

type requestProcessorHandler struct {
	connectionTable conntrack.ConnectionTable
}

func NewRequestProcessorHandler(table conntrack.ConnectionTable) srcds.Handler {
	return &requestProcessorHandler{
		connectionTable: table,
	}
}

func (h *requestProcessorHandler) Handle(responseWriter srcds.ConnectionWriter, msg srcds.Message, addr srcds.AddressPort) error {
	// Handle will handle the incoming connections to the proxy. It will forward every byte received to the proxy.
	// If it is a new connection, it will add an entry to the connection table and instantiate a controller that will
	// listen for responses from the proxy.
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
	} else {
		//log.Print("Known client: ", addr.String())
		serverConn, err = h.connectionTable.GetConnection(addr)
		if err != nil {
			return err
		}
	}
	return serverConn.Write(srcds.MessageToBytes(msg))
}
