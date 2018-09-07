package handler

import (
	"log"
	"srcds_proxy/proxy/srcds"
)

type LogHandler struct {
}

func (h *LogHandler) Handle(responseWriter srcds.ConnectionWriter, msg srcds.Message, addr srcds.AddressPort) error {
	log.Println("Received ", len(msg), " bytes from ", addr.String())
	return nil
}
