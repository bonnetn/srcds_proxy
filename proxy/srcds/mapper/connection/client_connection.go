package connection

import (
	"net"

	m "github.com/bonnetn/srcds_proxy/proxy/srcds/model"
	"github.com/bonnetn/srcds_proxy/utils"
	"github.com/golang/glog"
)

// ToClientConnectionWrapper creates a connection that uses a listening socket. You have to provide the address
// where to respond, because a listening connection is not connected to a specific host. You also have to provide
// the received packet in the MsgChan.
func ToClientConnectionWrapper(done <-chan utils.DoneEvent, conn *net.UDPConn, raddr net.UDPAddr) *m.ConnectionWrapper {

	outputChan := make(chan m.Message)
	inputChan := make(chan m.Message)
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
					glog.V(4).Info("Writing ", len(msg), " bytes to client.")
					if _, err := conn.WriteToUDP(msg, &raddr); err != nil {
						glog.Error("Error while writing to client: ", err)
						return
					}
					glog.V(4).Info("Successfully sent ", len(msg), " bytes to client.")
				}()
			}
		}
	}()

	return &m.ConnectionWrapper{
		MsgChan: inputChan,
		Connection: &connection{
			inputChannel:  inputChan,
			outputChannel: outputChan,
		},
	}
}
