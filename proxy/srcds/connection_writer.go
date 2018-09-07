package srcds

import (
	"net"
)

type ConnectionWriter interface {
	Write(Message) error
}

func NewConnectionWriter(conn net.UDPConn, raddr *net.UDPAddr) ConnectionWriter {
	// raddr can be nil in the case of a pre-connected connection.
	return &responseWriterImpl{
		conn:  conn,
		raddr: raddr,
	}
}

type responseWriterImpl struct {
	conn  net.UDPConn
	raddr *net.UDPAddr
}

func (c *responseWriterImpl) Write(msg Message) error {
	var err error
	if c.raddr != nil {
		_, err = c.conn.WriteToUDP(MessageToBytes(msg), c.raddr)
	} else {
		_, err = c.conn.Write(MessageToBytes(msg))
	}
	return err
}
