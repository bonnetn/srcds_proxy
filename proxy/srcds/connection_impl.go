package srcds

type connection struct {
	inputChannel  <-chan Message
	outputChannel chan<- Message
}

func (c *connection) OutputChannel() chan<- Message {
	return c.outputChannel
}

func (c *connection) InputChannel() <-chan Message {
	return c.inputChannel
}
