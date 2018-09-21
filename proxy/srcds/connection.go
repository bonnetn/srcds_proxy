package srcds

import "srcds_proxy/proxy/srcds/model"

type ConnectionWithPacketChan struct {
	MsgChan    chan model.Message
	Connection Connection
}

type Connection interface {
	InputChannel() <-chan model.Message
	OutputChannel() chan<- model.Message
}
