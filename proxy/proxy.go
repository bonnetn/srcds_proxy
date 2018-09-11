package proxy

import (
	"log"
	"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/config"
	)

func Launch() error {
	done := make(chan srcds.DoneEvent)
	defer close(done)

	log.Println("INFO: Starting proxy...")
	log.Println("INFO: Listening connections.")
	listener, err := srcds.Listen(done, config.ListenFullAddr)
	if err != nil {
		log.Println("ERR: Could not listen: ", err)
		return err
	}

	log.Println("INFO: Accepting connections.")
	bindings := srcds.AssociateWithServerConnection(done, listener.Accept(done))
	for bind := range bindings {
		forwardMessages(done, bind.ServerConnection, bind.ClientConnection)
		forwardMessages(done, bind.ClientConnection, bind.ServerConnection)
	}
	log.Println("INFO: Proxy stopped.")

	return nil
}

func forwardMessages(done <-chan srcds.DoneEvent, from, to srcds.Connection) {
	go func() {
		for msg := range from.InputChannel() {
			if srcds.IsDone(done) {
				return
			}
			to.OutputChannel() <- msg
		}
	}()
}