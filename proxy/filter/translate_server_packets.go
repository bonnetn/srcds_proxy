package filter

import (
	"bytes"

	"net"

	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

func translateSingleSvPacket(ctx models.ProxyContext, packet *models.Packet, clientConn *net.UDPConn) {
	if !bytes.Equal(packet.Src.IP[:], ctx.ServerHost.IP[:]) {
		return // Sent from a client, not from the server.
	}

	dst, ok := ctx.ServerToClientTbl[packet.Dst]
	if !ok {
		glog.Warningf("Could not route  response from server to %v", packet.Dst)
		return
	}

	packet.Socket = clientConn
	packet.Dst = *dst

}

// TranslateServerPackets takes packets sent from the server to the proxy and translates the destination address and connection to the clients.
func TranslateServerPackets(ctx models.ProxyContext, packetQueue <-chan models.Packet, clientConn *net.UDPConn) <-chan models.Packet {
	result := make(models.PacketQueue)
	go func() {
		for {
			pkt := <-packetQueue
			translateSingleSvPacket(ctx, &pkt, clientConn)
			result <- pkt
		}
	}()
	return result
}
