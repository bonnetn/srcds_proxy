package srcds

import (
	"net"
		"context"
	)

func Dial(addr string) (*net.UDPConn, error) {
	// Dial will create an UDP connection.
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	connection, err := net.DialUDP("udp", nil, laddr)
	if err != nil {
		return nil, err
	}

	if err = setSRCSConnectionOptions(connection); err != nil {
		return nil, err
	}

	return connection, nil
}

func Listen(ctx context.Context, addr string) (*net.UDPConn, error) {
	// Listen will create a listening UDP connection.
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	connection, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	// Ensure connection is always eventually closed.
	go func() {
		<-ctx.Done()
		connection.Close()
	}()

	if err = setSRCSConnectionOptions(connection); err != nil {
		return nil, err
	}

	return connection, nil
}

func Serve(ctx context.Context, connection net.UDPConn, handler Handler) error {
	// Serve will read data from a the connection to a buffer and call the handler provided.
	var (
		n          int
		sourceAddr *net.UDPAddr
		err        error
		buf        = make([]byte, MaxDatagramSize)
	)
	for {
		select {
			case <-ctx.Done():
				return nil
			default:
		}

		n, sourceAddr, err = connection.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		if n != 0 {
			msg := BytesToMessage(buf[:n])
			responseWriter := NewConnectionWriter(connection, sourceAddr) // object to respond
			if err = handler.Handle(responseWriter, msg, UDPAddrToAddressPort(*sourceAddr)); err != nil {
				return err
			}
		}
	}
}

func setSRCSConnectionOptions(connection *net.UDPConn) error {
	// Set the buffers size
	if err := connection.SetWriteBuffer(MaxDatagramSize); err != nil {
		return err
	}
	if err := connection.SetReadBuffer(MaxDatagramSize); err != nil {
		return err
	}
	return nil
}
