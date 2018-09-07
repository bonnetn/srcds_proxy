package srcds

import "context"

type Handler interface {
	Handle(context.Context, ConnectionWriter, Message, AddressPort) error
}
