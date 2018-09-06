package worker

import (
	"net"
	"srcds_proxy/proxy/model/srcds_connection"
	"srcds_proxy/proxy/controller"
)

type Worker interface {
	Join() error
}

type proxyWorker struct {
	conn     srcds_connection.SRCDSConnection
	handler  controller.Handler
	stopChan chan error
}

func NewProxyWorker(conn srcds_connection.SRCDSConnection, handler controller.Handler) Worker {
	pw := proxyWorker{
		conn:    conn,
		handler: handler,
	}
	go pw.run()
	return &pw
}

func (w *proxyWorker) Join() error {
	return <-w.stopChan
}

func (w *proxyWorker) run() {
	var (
		n    int
		addr *net.UDPAddr
		err  error
		buf  = make([]byte, srcds_connection.MaxDatagramSize)
	)

	for {
		n, addr, err = w.conn.ReadFromUDP(buf)
		if err != nil {
			w.stop(err)
		}

		if n != 0 {
			if err = w.handler.Handle(buf, *addr, n); err != nil {
				w.stop(err)
			}
		}
	}
	w.stop(nil)
}
func (w *proxyWorker) stop(err error) {
	w.stopChan <- err
}
