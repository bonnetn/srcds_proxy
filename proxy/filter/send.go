package filter

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

// SendQueue sends all the packets from the queue.
func SendQueue(packetQueue <-chan models.Packet, conn *net.UDPConn) {
	for {
		pkt := <-packetQueue

		var err error
		if *pkt.Socket == *conn {
			dst := models.HostToUDPAddr(&pkt.Dst)
			_, err = pkt.Socket.WriteToUDP(pkt.Content[:pkt.Size], dst)
		} else {
			_, err = pkt.Socket.Write(pkt.Content[:pkt.Size])
		}
		if err != nil {
			glog.Error("Could not send packet.", err)
		}
	}
}
