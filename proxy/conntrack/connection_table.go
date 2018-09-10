package conntrack

import (
	"errors"
	"sync"
	"srcds_proxy/proxy/srcds"
)

var NoConnectionAssociatedError = errors.New("no connection associated with that address")

type ConnectionTable interface {
	GetConnection(srcds.AddressPort) (srcds.ConnectionWriter, error)
	GetOrStoreConnection(srcds.AddressPort, srcds.ConnectionWriter) srcds.ConnectionWriter
	HasConnection(srcds.AddressPort) bool
	RemoveConnection(srcds.AddressPort)
}

var instance *connectionTableImpl
var once sync.Once

func GetConnectionTable() ConnectionTable {
	once.Do(func() {
		instance = &connectionTableImpl{
			Map: sync.Map{},
		}
	})
	return instance
}

type connectionTableImpl struct {
	sync.Map
}

func (tbl *connectionTableImpl) HasConnection(addr srcds.AddressPort) bool {
	_, ok := tbl.Load(addr)
	return ok
}

func (tbl *connectionTableImpl) GetConnection(addr srcds.AddressPort) (srcds.ConnectionWriter, error) {
	if conn, ok := tbl.Load(addr); ok {
		return conn.(srcds.ConnectionWriter), nil
	}
	return nil, NoConnectionAssociatedError
}

func (tbl *connectionTableImpl) GetOrStoreConnection(addr srcds.AddressPort, writer srcds.ConnectionWriter) srcds.ConnectionWriter {
	conn, _ := tbl.LoadOrStore(addr, writer)
	return conn.(srcds.ConnectionWriter)
}

func (tbl *connectionTableImpl) RemoveConnection(addr srcds.AddressPort) {
	tbl.Delete(addr)
}
