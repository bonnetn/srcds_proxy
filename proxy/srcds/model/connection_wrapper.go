package model

type ConnectionWrapper struct {
	// ConnectionWrapper is a wrapper that allows you to manually add messages in the InputChannel.
	MsgChan    chan Message
	Connection Connection
}
