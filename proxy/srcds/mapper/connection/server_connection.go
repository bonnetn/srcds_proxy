package connection

import (
	"net"

	m "github.com/bonnetn/srcds_proxy/proxy/srcds/model"
	"github.com/bonnetn/srcds_proxy/utils"
	"github.com/golang/glog"
)

// ToServerConnection created a connection that uses a dedicated socket to communicate with the server.
func ToServerConnection(done <-chan utils.DoneEvent, conn *net.UDPConn) m.Connection {
	// Listen on the connection and put all the messages received in the chan.
	inputChan := make(chan m.Message)
	outputChan := make(chan m.Message)

	go func() {
		defer close(inputChan)

		buffer := m.GetBufferPool().Get()
		defer m.GetBufferPool().Put(buffer)

		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if utils.IsDone(done) {
				return
			}
			if err != nil {
				glog.Error("Error while reading server response: ", err)
				return
			}
			inputChan <- m.BytesToMessage(buffer[:n])
		}
	}()

	go func() {
		defer close(outputChan)

		var msg m.Message
		for {
			select {
			case <-done:
				return
			case msg = <-outputChan:
				func() {
					defer m.GetBufferPool().Put(msg)
					glog.V(4).Info("Writing ", len(msg), " bytes to server.")
					if _, err := conn.Write(msg); err != nil {
						glog.Error("Error while writing to server: ", err)
						return
					}
					glog.V(4).Info("Successfully sent ", len(msg), " bytes to server.")
				}()
			}
		}
	}()
	return &connection{
		inputChannel:  inputChan,
		outputChannel: outputChan,
	}
}
