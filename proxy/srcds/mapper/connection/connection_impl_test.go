package connection

import (
	"testing"

	m "github.com/bonnetn/srcds_proxy/proxy/srcds/model"
)

func TestOutputChannel(t *testing.T) {
	var want chan m.Message
	c := &connection{
		inputChannel:  make(chan m.Message),
		outputChannel: want,
	}
	if got := c.OutputChannel(); got != want {
		t.Errorf("connection.OutputChannel() = %v, want %v", got, want)
	}
}

func TestInputChannel(t *testing.T) {
	var want chan m.Message
	c := &connection{
		outputChannel: make(chan m.Message),
		inputChannel:  want,
	}
	if got := c.InputChannel(); got != want {
		t.Errorf("connection.InputChannel() = %v, want %v", got, want)
	}
}
