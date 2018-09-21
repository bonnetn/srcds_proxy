package model

import (
	"testing"
)

func TestMessageToBytes(t *testing.T) {
	var (
		msg         = Message("test")
		want []byte = msg
		got         = MessageToBytes(msg)
	)
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("MessageToBytes() = %v, want %v", got, want)
			return
		}
	}
}

func TestBytesToMessageAssertBytesAreCopied(t *testing.T) {
	var (
		buf  = []byte("test")
		want = Message(buf)
		got  = BytesToMessage(buf)
	)
	got[0] = 'Z'
	if got[0] == want[0] {
		t.Errorf("MessageToBytes() = %v, want %v", got, want)
	}

}

func TestBytesToMessage(t *testing.T) {
	var (
		buf  = []byte("test")
		want = Message(buf)
		got  = BytesToMessage(buf)
	)
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("BytesToMessage() = %v, want %v", got, want)
			return
		}
	}
}
