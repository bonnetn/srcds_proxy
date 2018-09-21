package srcds

import (
	"net"
	"srcds_proxy/utils"
	"github.com/golang/glog"
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
				glog.V(1).Info("Packet received with no connection assigned, creating new connection.")
				result <- clientConn.Connection
				glog.V(1).Info("Connection created.")
			}
			msg := GetBufferPool().Get()
			copy(msg, buffer[:n])
			glog.V(3).Info("Received datagram of length ", n, " from a client.")
			clientConn.MsgChan <- msg[:n]
			glog.V(3).Info("Forwarded datagram of length ", n, " in the input channel.")
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
