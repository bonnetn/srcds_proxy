package srcds

type Handler interface {
	Handle(ConnectionWriter, Message, AddressPort) error
}
