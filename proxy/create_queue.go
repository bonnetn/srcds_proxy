package proxy

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/mapper"
	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

// createQueueFromConn creates a packet queue from a connection object, putting all the incoming packet into the queue.
func createQueueFromConn(listenAddr *models.Host) (models.PacketQueue, *net.UDPConn) {

	packetConn, err := net.ListenUDP("udp4", mapper.HostToUDPAddr(listenAddr))
	if err != nil {
		glog.Fatal(err)
	}

	packetQueue := make(models.PacketQueue)
	go func() {
		transferIncomingPackets(packetQueue, packetConn, *listenAddr)
	}()
	return packetQueue, packetConn

}

// TransferIncomingPackets listens on the given connection and put all the incoming packets in a queue.
func transferIncomingPackets(queue models.PacketQueue, conn *net.UDPConn, destination models.Host) {
	for {
		buf := make([]byte, models.UDPMaxSize)

		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			glog.Error("Could not read packet", err)
			continue
		}

		clientHost, err := mapper.UDPAddrToHost(clientAddr)
		if err != nil {
			glog.Error("Could not get client host.", err)
			continue
		}

		queue <- models.Packet{
			Src:     *clientHost,
			Dst:     destination,
			Size:    n,
			Content: buf,
			Socket:  conn,
		}
	}
}
