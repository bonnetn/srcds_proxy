package model

type ConnectionWithPacketChan struct {
	MsgChan    chan Message
	Connection Connection
}

type Connection interface {
	InputChannel() <-chan Message
	OutputChannel() chan<- Message
}
