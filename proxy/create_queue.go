package proxy

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

// createQueueFromConn creates a packet queue from a connection object, putting all the incoming packet into the queue.
func createQueueFromConn(listenAddr *models.Host) (models.PacketQueue, *net.UDPConn) {

	packetConn, err := net.ListenUDP("udp4", models.HostToUDPAddr(listenAddr))
	if err != nil {
		glog.Fatal(err)
	}

	packetQueue := make(models.PacketQueue)
	go func() {
		packetQueue.TransferIncomingPackets(packetConn, *listenAddr)
	}()
	return packetQueue, packetConn

}
