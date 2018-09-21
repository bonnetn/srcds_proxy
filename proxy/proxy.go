package proxy

import (
		"srcds_proxy/proxy/srcds"
	"srcds_proxy/proxy/config"
	"srcds_proxy/utils"
	"github.com/golang/glog"
	"srcds_proxy/proxy/srcds/model"
)

func Launch() error {
	done := make(chan utils.DoneEvent)
	defer close(done)

	glog.Info("Starting proxy.")

	glog.Info("Listening for new connections.")
	listener, err := srcds.Listen(done, config.ListenAddr())
	if err != nil {
		glog.Error("Could not listen", err)
		return err
	}

	glog.Info("Accepting connections.")
	bindings := srcds.AssociateWithServerConnection(done, listener.Accept(done))
	for bind := range bindings {
		glog.V(1).Info("New binding received, creating forward goroutines.")
		go forwardMessages(done, bind.ServerConnection, bind.ClientConnection)
		go forwardMessages(done, bind.ClientConnection, bind.ServerConnection)
	}
	glog.Info("Proxy stopped.")

	return nil
}

func forwardMessages(done <-chan utils.DoneEvent, from, to srcds.Connection) {
	var msg model.Message
	for {
		select {
		case <-done:
			return
		case msg = <-from.InputChannel():
			if len(msg) <= 0 {
				return
			}
			glog.V(2).Info("Forwarding a message of length ", len(msg), ".")
			to.OutputChannel() <- msg
			glog.V(2).Info("Successfully forwarded message of length ", len(msg),".")
		}
	}
}
