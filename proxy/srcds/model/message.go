package model

// Message is a SRCDS message that is exchanged between the SRCDS server and the clients.
type Message []byte

// MessageToBytes casts the message to raw bytes.
func MessageToBytes(m Message) []byte {
	return m
}

// BytesToMessage creates a new message that contains a copy of the buffer provided.
func BytesToMessage(b []byte) Message {
	c := GetBufferPool().Get()[:len(b)]
	copy(c, b)
	return c
}
