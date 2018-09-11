package srcds

const MaxDatagramSize = 4096

type Message []byte

func MessageToBytes(m Message) []byte {
	return m
}

func BytesToMessage(b []byte) Message {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
