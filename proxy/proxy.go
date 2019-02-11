package proxy

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/config"
	"github.com/bonnetn/srcds_proxy/proxy/filter"
	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

func makeInitialPacketQueue(listenAddr *models.Host) (models.PacketQueue, *net.UDPConn) {

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

// Launch launches the proxy.
func Launch() error {

	glog.Info("Starting proxy.")
	glog.Info("Listen address: ", config.ListenAddr())
	glog.Info("Proxy to address: ", config.ServerAddr())

	listenAddr, err := net.ResolveUDPAddr("udp4", config.ListenAddr())
	if err != nil {
		glog.Fatal(err)
	}

	listenHost, err := models.UDPAddrToHost(listenAddr)
	if err != nil {
		glog.Fatal(err)
	}

	dstAddr, err := net.ResolveUDPAddr("udp4", config.ServerAddr())
	if err != nil {
		glog.Fatal(err)
	}
	dstHost, err := models.UDPAddrToHost(dstAddr)
	if err != nil {
		glog.Fatal(err)
	}

	rootQueue, clientConn := makeInitialPacketQueue(listenHost)
	ctx := models.ProxyContext{
		ClientToServerTbl: &models.NatTable{},
		ServerToClientTbl: &models.NatTable{},
		ServerHost:        dstHost,
		ProxyHost:         listenHost,
		RootQueue:         rootQueue,
	}

	queue := (<-chan models.Packet)(ctx.RootQueue)
	queue = filter.TranslateClientPackets(ctx, queue, clientConn)
	queue = filter.TranslateServerPackets(ctx, queue, clientConn)
	filter.SendQueue(queue, clientConn)

	return nil

}
