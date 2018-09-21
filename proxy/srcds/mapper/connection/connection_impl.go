package connection

import m "github.com/bonnetn/srcds_proxy/proxy/srcds/model"

type connection struct {
	inputChannel  <-chan m.Message
	outputChannel chan<- m.Message
}

func (c *connection) OutputChannel() chan<- m.Message {
	return c.outputChannel
}

func (c *connection) InputChannel() <-chan m.Message {
	return c.inputChannel
}
