package model

// ConnectionWrapper is a wrapper that allows you to manually add messages in the InputChannel.
type ConnectionWrapper struct {
	MsgChan    chan Message
	Connection Connection
}
