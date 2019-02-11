package models

import (
	"net"

	"github.com/golang/glog"
)

// PacketQueue is a queue of packets.
type PacketQueue chan Packet

// TransferIncomingPackets listens on the given connection and put all the incoming packets in a queue.
func (queue PacketQueue) TransferIncomingPackets(conn *net.UDPConn, destination Host) {
	for {
		buf := make([]byte, UDPMaxSize)

		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			glog.Error("Could not read packet", err)
			continue
		}

		clientHost, err := UDPAddrToHost(clientAddr)
		if err != nil {
			glog.Error("Could not get client host.", err)
			continue
		}

		queue <- Packet{
			Src:     *clientHost,
			Dst:     destination,
			Size:    n,
			Content: buf,
			Socket:  conn,
		}
	}

}
