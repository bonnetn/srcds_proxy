package srcds

import "srcds_proxy/proxy/srcds/model"

type connection struct {
	inputChannel  <-chan model.Message
	outputChannel chan<- model.Message
}

func (c *connection) OutputChannel() chan<- model.Message {
	return c.outputChannel
}

func (c *connection) InputChannel() <-chan model.Message {
	return c.inputChannel
}
