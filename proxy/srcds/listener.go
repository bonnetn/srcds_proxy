package srcds

import (
	"net"
	"srcds_proxy/utils"
	"log"
)

type Listener struct {
	conn *net.UDPConn
}

var clientConnTable ConnectionTable

func (l *Listener) Accept(done chan utils.DoneEvent) <-chan Connection {
	result := make(chan Connection)
	go func() {
		defer close(result)

		buffer := GetBufferPool().Get()
		defer GetBufferPool().Put(buffer)

		for {
			n, raddr, err := l.conn.ReadFromUDP(buffer)
			if utils.IsDone(done) {
				return
			}
			if err != nil {
				return
			}


			clientConn, loaded := l.getOrCreateClientConn(done, raddr)
			if !loaded {
				result <- clientConn.Connection
				log.Println("DEBUG: New connection created.")
			}
			msg := GetBufferPool().Get()
			copy(msg, buffer[:n])
			log.Println("DEBUG: New datagram received from world.")
			clientConn.MsgChan <- msg[:n]
		}

	}()
	return result
}

func (l *Listener) getOrCreateClientConn(done <-chan utils.DoneEvent, raddr *net.UDPAddr) (*ConnectionWithPacketChan, bool) {
	// Create a new connection.
	killNewClientConn := make(chan utils.DoneEvent)
	newClientConn := NewConnectionWithPacketChan(channelOr(done, killNewClientConn), l.conn, *raddr)

	clientConn, loaded := clientConnTable.GetOrReplace(UDPAddrToAddressPort(*raddr), newClientConn)
	if loaded {
		close(killNewClientConn) // If this connection is not used, kill the workers related to that connection.
	}
	return clientConn, loaded
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
