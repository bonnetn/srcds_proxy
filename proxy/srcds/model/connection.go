package model

type Connection interface {
	InputChannel() <-chan Message
	OutputChannel() chan<- Message
}
