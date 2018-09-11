package srcds

import (
	"net"
	"srcds_proxy/utils"
)

type Listener struct {
	conn *net.UDPConn
}

var clientConnTable ConnectionTable

func (l *Listener) Accept(done chan utils.DoneEvent) <-chan Connection {
	result := make(chan Connection)
	go func() {
		defer close(result)

		for {
			buffer := make([]byte, MaxDatagramSize)
			n, raddr, err := l.conn.ReadFromUDP(buffer)
			if utils.IsDone(done) {
				return
			}
			if err != nil {
				return
			}
			addr := UDPAddrToAddressPort(*raddr)

			killChan := make(chan utils.DoneEvent)
			clientConnection := NewClientConnection(channelOr(done, killChan), l.conn, *raddr, BytesToMessage(buffer[:n]))
			conn, loaded := clientConnTable.GetOrReplace(addr, clientConnection)
			if !loaded {
				result <- conn
			} else {
				close(killChan) // If this connection is not used, kill all the workers.
			}
		}

	}()
	return result
}

func channelOr(a, b <-chan utils.DoneEvent) <-chan utils.DoneEvent {
	orChan := make(chan utils.DoneEvent)
	go func() {
		defer close(orChan)
		for {
			select {
			case <-a:
				return
			case <-b:
				return
			}
		}
	}()
	return orChan
}
