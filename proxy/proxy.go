package proxy

import (
	"github.com/bonnetn/srcds_proxy/proxy/config"
	"github.com/bonnetn/srcds_proxy/proxy/filter"
	"github.com/bonnetn/srcds_proxy/proxy/mapper"
	"github.com/bonnetn/srcds_proxy/proxy/models"
	"github.com/golang/glog"
)

// Launch launches the proxy.
func Launch() error {

	glog.Info("Starting proxy.")
	glog.Info("Listen address: ", config.ListenAddr())
	glog.Info("Proxy to address: ", config.ServerAddr())

	listenHost, err := mapper.StringToHost(config.ListenAddr())
	if err != nil {
		glog.Fatal(err)
	}

	dstHost, err := mapper.StringToHost(config.ServerAddr())
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
