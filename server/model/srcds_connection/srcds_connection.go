package srcds_connection

import (
	"net"
)

const MaxDatagramSize = 4096

type SRCDSConnection struct {
	net.UDPConn
}

func NewSRCDSOutboundConnection(outAddr string) (*SRCDSConnection, error) {
	addr, err := net.ResolveUDPAddr("udp", outAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(MaxDatagramSize)
	conn.SetWriteBuffer(MaxDatagramSize)
	return &SRCDSConnection{
		UDPConn: *conn,
	}, nil
}

func NewSRCDSInboundConnection(listenAddr string) (*SRCDSConnection, error) {
	addr, err := net.ResolveUDPAddr("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	conn.SetReadBuffer(MaxDatagramSize)
	conn.SetWriteBuffer(MaxDatagramSize)
	return &SRCDSConnection{
		UDPConn: *conn,
	}, nil
}
