package models

// ProxyContext is the context of the proxy while applying filters to the packets.
type ProxyContext struct {
	ClientToServerTbl *NatTable // Client IP/Port --> Connection  to the server
	ServerToClientTbl *NatTable // Local Port --> Client IP/Port

	ProxyHost  *Host
	ServerHost *Host

	RootQueue PacketQueue
}
