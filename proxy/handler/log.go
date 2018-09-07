package handler

import (
	"log"
	"srcds_proxy/proxy/srcds"
	"context"
)

type LogHandler struct {
}

func (h *LogHandler) Handle(_ context.Context, responseWriter srcds.ConnectionWriter, msg srcds.Message, addr srcds.AddressPort) error {
	log.Println("INFO: Received ", len(msg), " bytes from ", addr.String())
	return nil
}
