package worker

import (
	"net"
	"errors"
	"srcds_proxy/server/handler"
)

type udpListener struct {
	conn     *net.UDPConn
	handler  handler.Handler
	stopChan chan struct{}
}

func NewUDPListenerWorker(conn *net.UDPConn, handler handler.Handler) (*udpListener, error) {
	if conn == nil {
		return nil, errors.New("nil connection")
	}

	return &udpListener{
		conn:    conn,
		handler: handler,
	}, nil

}

func (w *udpListener) Run() error {
	var (
		n    int
		addr *net.UDPAddr
		err  error
		buf  = make([]byte, MAX_DATAGRAM_SIZE)
	)

	defer w.stop()
	for {
		n, addr, err = w.conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		if n != 0 {
			if err = w.handler.Handle(buf, addr, n); err != nil {
				return err
			}
		}
	}
}

func (w *udpListener) Join() {
	<-w.stopChan
}

func (w *udpListener) stop() {
	w.stopChan <- struct{}{}
}
