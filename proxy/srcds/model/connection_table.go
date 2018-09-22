package model

import (
	"sync"
)

// ConnectionTable keeps track of the binding between a client and its connection to the server.
// It is thread-safe.
type ConnectionTable struct {
	sync.Map
}

// GetOrReplace returns the value of an entry in the table. If that entry does not exist, it creates one.
func (tbl *ConnectionTable) GetOrReplace(addr AddressPort, conn *ConnectionWrapper) (*ConnectionWrapper, bool) {
	res, loaded := tbl.LoadOrStore(addr, conn)
	return res.(*ConnectionWrapper), loaded
}

// Remove removes an entry from the ConnectionTable.
func (tbl *ConnectionTable) Remove(addr AddressPort) {
	tbl.Delete(addr)
}
