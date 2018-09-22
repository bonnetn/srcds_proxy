package model

// Connection is an interface that allows to send and receive messages via channels.
type Connection interface {
	InputChannel() <-chan Message
	OutputChannel() chan<- Message
}
