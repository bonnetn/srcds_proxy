package srcds

import (
	"net"
	"srcds_proxy/utils"
	"log"
)

type ConnectionWithPacketChan struct {
	MsgChan    chan Message
	Connection Connection
}

type Connection interface {
	InputChannel() <-chan Message
	OutputChannel() chan<- Message
}

func NewConnectionWithPacketChan(done <-chan utils.DoneEvent, conn *net.UDPConn, raddr net.UDPAddr) *ConnectionWithPacketChan {
	// NewConnectionWithPacketChan creates a connection that uses a listening socket. You have to provide the address
	// where to respond, because a listening connection is not connected to a specific host. You also have to provide
	// the received packet in the MsgChan.

	outputChan := make(chan Message)
	inputChan := make(chan Message)
	go func() {
		defer close(outputChan)

		var msg Message
		for {
			select {
			case <-done:
				return
			case msg = <-outputChan:
				log.Println("DEBUG: [CLIENT_CON] Writing to client", len(msg))
				conn.WriteToUDP(msg, &raddr)
				log.Println("DEBUG: [CLIENT_CON] Wrote to client", len(msg))
				GetBufferPool().Put(msg)
				log.Println("DEBUG: [CLIENT_CON] freed buffer", len(msg))
			}
		}
	}()

	return &ConnectionWithPacketChan{
		MsgChan: inputChan,
		Connection: &connection{
			inputChannel:  inputChan,
			outputChannel: outputChan,
		},
	}
}

func NewConnection(done <-chan utils.DoneEvent, conn *net.UDPConn) Connection {
	// NewConnection created a connection that uses a dedicated socket to communicate with the server.

	// Listen on the connection and put all the messages received in the chan.
	inputChan := make(chan Message)
	outputChan := make(chan Message)

	go func() {
		defer close(inputChan)

		buffer := GetBufferPool().Get()
		defer GetBufferPool().Put(buffer)

		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if utils.IsDone(done) {
				return
			}
			if err != nil {
				log.Println("ERROR: Error while reading server response.", err)
				return
			}
			inputChan <- BytesToMessage(buffer[:n])
		}
	}()

	go func() {
		defer close(outputChan)

		var msg Message
		for {
			select {
			case <-done:
				return
			case msg = <-outputChan:
				log.Println("DEBUG: [SERV_CON] Write to server", len(msg))
				conn.Write(msg)
				log.Println("DEBUG: [SERV_CON] Wrote to server", len(msg))
				GetBufferPool().Put(msg)
				log.Println("DEBUG: [SERV_CON] Freed buffer", len(msg))
			}
		}
	}()
	return &connection{
		inputChannel:  inputChan,
		outputChannel: outputChan,
	}
}

type connection struct {
	inputChannel  <-chan Message
	outputChannel chan<- Message
}

func (c *connection) OutputChannel() chan<- Message {
	return c.outputChannel
}

func (c *connection) InputChannel() <-chan Message {
	return c.inputChannel
}
