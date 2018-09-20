package proxy

import (
	"log"
	"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/config"
	"srcds_proxy/utils"
)

func Launch() error {
	done := make(chan utils.DoneEvent)
	defer close(done)

	log.Println("INFO: Starting proxy...")
	log.Println("INFO: Listening connections.")
	listener, err := srcds.Listen(done, config.ListenAddr())
	if err != nil {
		log.Println("ERR: Could not listen: ", err)
		return err
	}

	log.Println("INFO: Accepting connections.")
	bindings := srcds.AssociateWithServerConnection(done, listener.Accept(done))
	for bind := range bindings {
		go forwardMessages(done, bind.ServerConnection, bind.ClientConnection)
		go forwardMessages(done, bind.ClientConnection, bind.ServerConnection)
	}
	log.Println("INFO: Proxy stopped.")

	return nil
}

func forwardMessages(done <-chan utils.DoneEvent, from, to srcds.Connection) {
	var msg srcds.Message
	for {
		select {
		case <-done:
			return
		case msg = <-from.InputChannel():
			to.OutputChannel() <- msg
		}
	}
}
