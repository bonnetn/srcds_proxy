package srcds

const MaxDatagramSize = 4096

type Message []byte

func MessageToBytes(m Message) []byte {
	return m
}

func BytesToMessage(b []byte) Message {
	return b
}