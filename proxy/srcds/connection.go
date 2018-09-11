package srcds

import (
	"net"
	"srcds_proxy/utils"
)

type Connection interface {
	InputChannel() <-chan Message
	OutputChannel() chan<- Message
}

func NewClientConnection(done <-chan utils.DoneEvent, conn *net.UDPConn, raddr net.UDPAddr, initialMsg Message) Connection {
	// Client connection is a connection that uses a listening socket. You have to provide the address where to respond,
	// because a listening connection is not connected to a specific host.

	outputChan := make(chan Message)
	inputChan := make(chan Message)
	go func() {
		inputChan <- initialMsg
	}()
	go func() {
		defer close(outputChan)

		var msg Message
		for {
			select {
			case <-done:
				return
			case msg = <-outputChan:
				conn.WriteToUDP(MessageToBytes(msg), &raddr)
			}
		}
	}()

	return &connection{
		inputChannel:  inputChan,
		outputChannel: outputChan,
	}
}

func NewServerConnection(done <-chan utils.DoneEvent, conn *net.UDPConn) Connection {
	// Server connection is a connection that uses a dedicated socket to communicate with the server.

	// Listen on the connection and put all the messages recevied in the chan.
	inputChan := make(chan Message)
	outputChan := make(chan Message)
	go func() {
		defer close(inputChan)
		defer close(outputChan)

		var (
			msg    Message
			n      int
			err    error
			buffer = make([]byte, MaxDatagramSize)
		)
		for {
			select {
			case <-done:
				return
			case msg = <-outputChan:
				conn.Write(MessageToBytes(msg))
			default:
				n, _, err = conn.ReadFromUDP(buffer)
				if err != nil {
					return
				}
				inputChan <- BytesToMessage(buffer[:n])
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
