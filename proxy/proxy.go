package proxy

import (
	"log"
	"runtime"
	"srcds_proxy/proxy/worker"
	"srcds_proxy/proxy/model/conntrack"
	"srcds_proxy/proxy/controller/srcds"
	"srcds_proxy/proxy/model/srcds_connection"
)

const (
	ListenAddr = "" // Listen to every address
	ListenPort = "1234"
	ServerAddr = "127.0.0.1"
	ServerPort = "27016"

	ListenFullAddr = ListenAddr + ":" + ListenPort
	ServerFullAddr = ServerAddr + ":" + ServerPort
)

func Launch() error {

	listenConn, err := srcds_connection.NewSRCDSInboundConnection(ListenFullAddr)
	if err != nil {
		log.Print("could not create input connection: ", err)
		return err
	}
	defer listenConn.Close()

	connectionTable := conntrack.NewConnectionTable(ServerFullAddr)
	defer connectionTable.CloseAllConnections()

	var (
		workerCount = runtime.NumCPU()
		workers     = make([]worker.Worker, workerCount)
	)
	log.Println("Starting proxy workers...")
	for i := 0; i < workerCount; i++ {
		workers[i] = worker.NewProxyWorker(*listenConn, srcds.NewClientHandler(connectionTable, *listenConn))
	}

	log.Println("Proxy started! Listening for datagrams...")
	for i := 0; i < workerCount; i++ {
		if err = workers[i].Join(); err != nil {
			log.Println("Worker ", i, "crashed: ", err)
		}
	}
	log.Println("Proxy stopped.")

	return nil
}
