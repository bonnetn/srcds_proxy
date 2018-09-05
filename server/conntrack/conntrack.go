package conntrack

import (
	"net"
	"sync"
	"srcds_proxy/server/worker"
	"srcds_proxy/server/handler"
)

type ConntrackTable interface {
	GetOrCreateConnection(udpAddr net.UDPAddr) (*net.UDPConn, error)
	CloseAllConnections()
}

type conntrackTableImpl struct {
	sync.Mutex
	connectionTbl map[string]*net.UDPConn
	serverAddr    *net.UDPAddr
}

func NewConntrackTable(serverAddr net.UDPAddr) ConntrackTable {
	return &conntrackTableImpl{
		Mutex:         sync.Mutex{},
		connectionTbl: map[string]*net.UDPConn{},
		serverAddr:    &serverAddr,
	}
}

// TODO: make one get and one create to create the worker elsewhere

func (tbl *conntrackTableImpl) GetOrCreateConnection(udpAddr net.UDPAddr) (*net.UDPConn, error) {
	tbl.Lock()
	defer tbl.Unlock()

	// Return immediately if connection is already created.
	if con := tbl.connectionTbl[udpAddr.String()]; con != nil {
		return tbl.connectionTbl[udpAddr.String()], nil
	}

	// Create a connection with an available random port.
	conn, err := tbl.createConnection()
	if err != nil {
		return nil, err
	}

	// Create a worker that will process the responses from the server.
	_, err = worker.NewUDPListenerWorker(conn, handler.NewServerHandler(udpAddr))
	if err != nil {
		return nil, err
	}

	tbl.connectionTbl[udpAddr.String()] = conn
	return conn, nil
}

func (tbl *conntrackTableImpl) CloseAllConnections() {
	tbl.Lock()
	defer tbl.Unlock()

	for k, v := range tbl.connectionTbl {
		if v != nil {
			v.Close()
			delete(tbl.connectionTbl, k)
		}
	}
}

func (tbl *conntrackTableImpl) createConnection() (*net.UDPConn, error) {
	conn, err := net.ListenUDP("udp", tbl.serverAddr)
	if err != nil {
		return nil, err
	}
	conn.SetReadBuffer(worker.MAX_DATAGRAM_SIZE)
	conn.SetWriteBuffer(worker.MAX_DATAGRAM_SIZE)
	return conn, nil
}
