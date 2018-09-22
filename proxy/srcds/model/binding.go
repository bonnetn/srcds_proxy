package model

// Binding is a binding between the connection of the client to the proxy and the connection from the proxy to the
// server. This can be seen as: Client <=[ServConn]=> Proxy <=[ClientConn]=> ServerSRCDS
type Binding struct {
	ServerConnection Connection
	ClientConnection Connection
}
