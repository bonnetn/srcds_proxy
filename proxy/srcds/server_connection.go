package srcds

import (
	"srcds_proxy/utils"
	"net"
	"github.com/golang/glog"
)

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
				glog.Error("Error while reading server response: ", err)
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
				glog.V(4).Info("Writing ", len(msg), " bytes to server.")
				conn.Write(msg)
				glog.V(4).Info("Successfully sent ", len(msg), " bytes to server.")
				GetBufferPool().Put(msg)
				glog.V(4).Info("Freed the buffer.")
			}
		}
	}()
	return &connection{
		inputChannel:  inputChan,
		outputChannel: outputChan,
	}
}
