package srcds

import (
	"sync"
)

type ConnectionTable struct {
	sync.Map
}

func (tbl *ConnectionTable) GetOrReplace(addr AddressPort, conn *ConnectionWithPacketChan) (*ConnectionWithPacketChan, bool) {
	res, loaded := tbl.LoadOrStore(addr, conn)
	return res.(*ConnectionWithPacketChan), loaded
}

func (tbl *ConnectionTable) Remove(addr AddressPort) {
	tbl.Delete(addr)
}
