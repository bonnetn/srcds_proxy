package srcds

import (
	"net"
	"srcds_proxy/proxy/config"
	"srcds_proxy/utils"
	"log"
)

func Listen(done <-chan utils.DoneEvent, addr string) (*Listener, error) {
	conn, err := makeConnection(done, addr, true)
	if err != nil {
		return nil, err
	}
	return &Listener{
		conn: conn,
	}, err

}

func AssociateWithServerConnection(done <-chan utils.DoneEvent, connChan <-chan Connection) <-chan Binding {
	result := make(chan Binding)
	go func() {
		defer close(result)

		for clientConnection := range connChan {
			if utils.IsDone(done) {
				return
			}

			udpConn, err := dial(done, config.ServerAddr())
			if err != nil {
				return
			}
			log.Println("DEBUG: New server connection created.")

			result <- Binding{
				ServerConnection: NewConnection(done, udpConn),
				ClientConnection: clientConnection,
			}
		}
	}()
	return result
}

func dial(done <-chan utils.DoneEvent, addr string) (*net.UDPConn, error) {
	return makeConnection(done, addr, false)
}

func makeConnection(done <-chan utils.DoneEvent, addr string, listening bool) (*net.UDPConn, error) {

	// Listen will create a listening UDP ClientConnection.

	// First create the UDP listening ClientConnection.
	laddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	var udpConn *net.UDPConn
	if listening {
		udpConn, err = net.ListenUDP("udp", laddr)
	} else {
		udpConn, err = net.DialUDP("udp", nil, laddr)
	}
	if err != nil {
		return nil, err
	}

	if err = setSRCSConnectionOptions(udpConn); err != nil {
		return nil, err
	}

	// Close on done
	go func() {
		<-done
		udpConn.Close()
	}()

	return udpConn, nil
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
