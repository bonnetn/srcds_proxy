package server

import (
	"net"
	"log"
	"runtime"
	"srcds_proxy/server/worker"
	"srcds_proxy/server/conntrack"
	"srcds_proxy/server/handler"
)

const (
	ListenAddr = "" // Listen to every address
	ListenPort = "1234"
	ServerAddr = "localhost"
	ServerPort = "27015"

	ListenFullAddr = ListenAddr + ":" + ListenPort
	ServerFullAddr = ServerAddr + ":" + ServerPort
)

func Launch() error {

	listenConn, err := createListenConnection(ListenFullAddr)
	if err != nil {
		log.Print("could not create input connection: ", err)
		return err
	}
	defer listenConn.Close()

	serverAddr, err := net.ResolveUDPAddr("udp", ServerFullAddr)
	if err != nil {
		log.Print("could not resolve server address: ", err)
		return err
	}

	conntrack := conntrack.NewConntrackTable(*serverAddr)
	defer conntrack.CloseAllConnections()

	var (
		workerCount = runtime.NumCPU()
		workers     = make([]worker.Worker, workerCount)
	)
	for i := 0; i < workerCount; i++ {
		workers[i], err = worker.NewUDPListenerWorker(listenConn, handler.NewClientHandler(conntrack))
		if err != nil {
			log.Print("could not instantiate worker: ", err)
			return err
		}
		go launchWorker(workers[i])
	}

	log.Println("Proxy started! Listening for datagrams...")
	for i := 0; i < workerCount; i++ {
		workers[i].Join()
	}
	log.Println("Proxy stopped.")

	return nil
}

func launchWorker(worker worker.Worker) {
	if err := worker.Run(); err != nil {
		log.Print("worker crashed: ", err)
	}
}

func createListenConnection(listenAddr string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(worker.MAX_DATAGRAM_SIZE)
	conn.SetWriteBuffer(worker.MAX_DATAGRAM_SIZE)
	return conn, nil
}
