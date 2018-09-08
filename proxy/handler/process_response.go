package handler

import (
	"srcds_proxy/proxy/srcds"
	"context"
	"log"
)

type responseProcessorHandler struct {
	// responseProcessorHandler forwards the server responses to the client.
	clientWriter srcds.ConnectionWriter
}

func NewResponseProcessorHandler(clientWriter srcds.ConnectionWriter) srcds.Handler {
	return &responseProcessorHandler{
		clientWriter: clientWriter,
	}
}

func (h *responseProcessorHandler) Handle(
	// Handle forwards back every message to the clientWriter.
	ctx context.Context, responseWriter srcds.ConnectionWriter, msg srcds.Message, addr srcds.AddressPort) error {
	log.Print("received response ", len(msg), " bytes")
	h.clientWriter.Write(msg)
	return nil
}
