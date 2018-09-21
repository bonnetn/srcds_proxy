package srcds

import (
	"srcds_proxy/utils"
	"net"
	"github.com/golang/glog"
	"srcds_proxy/proxy/srcds/model"
)

func NewConnection(done <-chan utils.DoneEvent, conn *net.UDPConn) Connection {
	// NewConnection created a connection that uses a dedicated socket to communicate with the server.

	// Listen on the connection and put all the messages received in the chan.
	inputChan := make(chan model.Message)
	outputChan := make(chan model.Message)

	go func() {
		defer close(inputChan)

		buffer := model.GetBufferPool().Get()
		defer model.GetBufferPool().Put(buffer)

		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if utils.IsDone(done) {
				return
			}
			if err != nil {
				glog.Error("Error while reading server response: ", err)
				return
			}
			inputChan <- model.BytesToMessage(buffer[:n])
		}
	}()

	go func() {
		defer close(outputChan)

		var msg model.Message
		for {
			select {
			case <-done:
				return
			case msg = <-outputChan:
				glog.V(4).Info("Writing ", len(msg), " bytes to server.")
				conn.Write(msg)
				glog.V(4).Info("Successfully sent ", len(msg), " bytes to server.")
				model.GetBufferPool().Put(msg)
				glog.V(4).Info("Freed the buffer.")
			}
		}
	}()
	return &connection{
		inputChannel:  inputChan,
		outputChannel: outputChan,
	}
}
