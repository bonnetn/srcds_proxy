package model

import (
	"sync"
)

type ConnectionTable struct {
	sync.Map
}

func (tbl *ConnectionTable) GetOrReplace(addr AddressPort, conn *ConnectionWrapper) (*ConnectionWrapper, bool) {
	res, loaded := tbl.LoadOrStore(addr, conn)
	return res.(*ConnectionWrapper), loaded
}

func (tbl *ConnectionTable) Remove(addr AddressPort) {
	tbl.Delete(addr)
}
