package conntrack

import (
	"net"
	"sync"
	"errors"
	"srcds_proxy/proxy/model/srcds_connection"
)

var NoConnectionAssociatedError = errors.New("no connection associated with that address")

type ConnectionTable interface {
	GetConnection(udpAddr net.UDPAddr) (*srcds_connection.SRCDSConnection, error)
	CreateConnection(udpAddr net.UDPAddr) (*srcds_connection.SRCDSConnection, error)
	CloseAllConnections()
}

type connectionTableImpl struct {
	sync.Mutex
	connectionTbl map[string]srcds_connection.SRCDSConnection
	serverAddr    string
}

func NewConnectionTable(serverAddr string) (ConnectionTable) {
	// NewConnectionTable will create a connection table that maps client sockets to proxy sockets.

	return &connectionTableImpl{
		Mutex:         sync.Mutex{},
		connectionTbl: map[string]srcds_connection.SRCDSConnection{},
		serverAddr:    serverAddr,
	}
}

func (tbl *connectionTableImpl) GetConnection(udpAddr net.UDPAddr) (*srcds_connection.SRCDSConnection, error) {
	tbl.Lock()
	defer tbl.Unlock()

	if conn, ok := tbl.connectionTbl[udpAddr.String()]; ok {
		return &conn, nil
	}
	return nil, NoConnectionAssociatedError
}

func (tbl *connectionTableImpl) CreateConnection(udpAddr net.UDPAddr) (*srcds_connection.SRCDSConnection, error) {
	if _, ok := tbl.connectionTbl[udpAddr.String()]; ok {
		return nil, errors.New("there is already a connection associated with that address")
	}

	tbl.Lock()
	defer tbl.Unlock()

	// Create a connection with an available random port.
	conn, err := srcds_connection.NewSRCDSOutboundConnection(tbl.serverAddr)
	if err != nil {
		return nil, err
	}

	tbl.connectionTbl[udpAddr.String()] = *conn
	return conn, nil
}

func (tbl *connectionTableImpl) CloseAllConnections() {
	tbl.Lock()
	defer tbl.Unlock()

	for k, v := range tbl.connectionTbl {
		v.Close()
		delete(tbl.connectionTbl, k)
	}
}
