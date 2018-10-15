package model

import "sync"

// ConnectionTable keeps track of the binding between a client and its connection to the server.
type ConnectionTable interface {
	GetOrReplace(addr AddressPort, conn *ConnectionWrapper) (*ConnectionWrapper, bool)
	Remove(addr AddressPort)
}

// connectionTable is a thread-safe implementation of the connection table which stores connections locally.
type connectionTable struct {
	sync.Map
}

// NewConnectionTable creates a local connection table.
func NewConnectionTable() ConnectionTable {
	return &connectionTable{}
}

// GetOrReplace returns the value of an entry in the table. If that entry does not exist, it creates one.
func (tbl *connectionTable) GetOrReplace(addr AddressPort, conn *ConnectionWrapper) (*ConnectionWrapper, bool) {
	res, loaded := tbl.LoadOrStore(addr, conn)
	return res.(*ConnectionWrapper), loaded
}

// Remove removes an entry from the ConnectionTable.
func (tbl *connectionTable) Remove(addr AddressPort) {
	tbl.Delete(addr)
}
