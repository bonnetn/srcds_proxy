package proxy

import (
	"net"

	"github.com/bonnetn/srcds_proxy/proxy/config"
	"github.com/bonnetn/srcds_proxy/proxy/filter"
	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

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

	rootQueue, clientConn := createQueueFromConn(listenHost)
	ctx := models.ProxyContext{
		ClientToServerTbl: &models.NatTable{},
		ServerToClientTbl: map[models.Host]*models.Host{},
		ServerHost:        dstHost,
		ProxyHost:         listenHost,
		RootQueue:         rootQueue,
	}

	queue := (<-chan models.Packet)(ctx.RootQueue)
	queue = filter.TranslateClientPackets(ctx, queue)
	queue = filter.TranslateServerPackets(ctx, queue, clientConn)
	filter.SendQueue(queue, clientConn)

	return nil

}
