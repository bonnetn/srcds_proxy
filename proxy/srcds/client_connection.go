package srcds

import (
	"srcds_proxy/utils"
	"net"
	"github.com/golang/glog"
	m "srcds_proxy/proxy/srcds/model"
)

func NewConnectionWithPacketChan(done <-chan utils.DoneEvent, conn *net.UDPConn, raddr net.UDPAddr) *m.ConnectionWithPacketChan {
	// NewConnectionWithPacketChan creates a connection that uses a listening socket. You have to provide the address
	// where to respond, because a listening connection is not connected to a specific host. You also have to provide
	// the received packet in the MsgChan.

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
				glog.V(4).Info("Writing ", len(msg), " bytes to client.")
				conn.WriteToUDP(msg, &raddr)
				glog.V(4).Info("Successfully sent ", len(msg), " bytes to client.")
				m.GetBufferPool().Put(msg)
				glog.V(4).Info("Freed the buffer.")
			}
		}
	}()

	return &m.ConnectionWithPacketChan{
		MsgChan: inputChan,
		Connection: &connection{
			inputChannel:  inputChan,
			outputChannel: outputChan,
		},
	}
}
