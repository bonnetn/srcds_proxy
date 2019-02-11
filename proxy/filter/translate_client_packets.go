package filter

import (
	"bytes"
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

func translateSingleClPacket(ctx models.ProxyContext, packet *models.Packet) error {
	if bytes.Equal(packet.Src.IP[:], ctx.ServerHost.IP[:]) {
		return nil // Sent from the server, not from a client.
	}

	conn, ok := ctx.ClientToServerTbl.LoadConnection(packet.Src)
	if ok {
		// Client connection is already known by the proxy. We can find the server connection associated.
		// Just translate the address and return.
		packet.Socket = conn
		packet.Dst = *ctx.ServerHost
		return nil
	}

	// Proxy does not know this client connection. It will create a server connection and populate the connection table.
	glog.Infof("New connection received from %v.", packet.Src)

	// Create connection
	conn, err := net.DialUDP("udp4", nil, models.HostToUDPAddr(ctx.ServerHost))
	if err != nil {
		return err
	}

	// Try to store the connection
	newConn, loaded := ctx.ClientToServerTbl.LoadOrStoreConnection(packet.Src, conn)
	if loaded {
		// If there is already a connection, there is a data race.
		// Close the (unused) new connection and translate the packet with the connection
		// present in the table.
		err = conn.Close()
		packet.Socket = newConn
		packet.Dst = models.Host(*ctx.ServerHost)
		return err
	}

	// Also store the conn in the Server -> Client table.
	addr, err := net.ResolveUDPAddr(conn.LocalAddr().Network(), conn.LocalAddr().String())
	if err != nil {
		return err
	}

	// Try to store the binding
	localHost, err := models.UDPAddrToHost(addr)
	if err != nil {
		return err
	}

	ctx.ServerToClientTbl[*localHost] = &packet.Src

	// Make a new worker that will put  the incoming packets into the queue.
	if err := createWorker(ctx, newConn); err != nil {
		return err
	}

	packet.Socket = newConn
	packet.Dst = models.Host(*ctx.ServerHost)
	return nil

}

func createWorker(ctx models.ProxyContext, conn *net.UDPConn) error {
	destAddr, err := net.ResolveUDPAddr(conn.LocalAddr().Network(), conn.LocalAddr().String())
	if err != nil {
		return err
	}
	destHost, err := models.UDPAddrToHost(destAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			buf := make([]byte, models.UDPMaxSize)

			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				glog.Error("Could not read packet", err)
				continue
			}

			clientHost, err := models.UDPAddrToHost(clientAddr)
			if err != nil {
				glog.Error("Could not get client host.", err)
				continue
			}

			ctx.RootQueue <- models.Packet{
				Src:     *clientHost,
				Dst:     *destHost,
				Size:    n,
				Content: buf,
				Socket:  conn,
			}
		}
	}()
	return nil
}

// TranslateClientPackets takes packets sent to the proxy and translate the IP and connection to the server.
func TranslateClientPackets(ctx models.ProxyContext, packetQueue <-chan models.Packet) <-chan models.Packet {
	result := make(models.PacketQueue)
	go func() {
		for {
			pkt := <-packetQueue
			err := translateSingleClPacket(ctx, &pkt)
			if err != nil {
				glog.Error("Could not translate client packet.", err)
			}
			result <- pkt
		}
	}()
	return result
}
