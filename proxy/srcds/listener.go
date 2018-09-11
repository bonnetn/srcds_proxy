package srcds

import (
	"net"
)

type Listener struct {
	conn *net.UDPConn
}

var clientConnTable ConnectionTable

func (l *Listener) Accept(done chan DoneEvent) <-chan Connection {
	result := make(chan Connection)
	go func() {
		defer close(result)

		for {
			buffer := make([]byte, MaxDatagramSize)
			n, raddr, err := l.conn.ReadFromUDP(buffer)
			if IsDone(done) {
				return
			}
			if err != nil {
				return
			}

			clientConnection := NewClientConnection(done, l.conn, *raddr, BytesToMessage(buffer[:n]))
			conn, loaded := clientConnTable.GetOrReplace(UDPAddrToAddressPort(*raddr), clientConnection)
			if !loaded {
				result <- conn
			}
		}

	}()
	return result
}
