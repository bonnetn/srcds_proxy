package model

import (
	"sync"
	"github.com/bonnetn/srcds_proxy/proxy/config"
)

type bufferPool struct {
	pool sync.Pool
}

var (
	singletonBufferPool *bufferPool
	once                sync.Once
)

func GetBufferPool() *bufferPool {
	once.Do(func() {
		singletonBufferPool = &bufferPool{
			pool: sync.Pool{
				New: newBuffer,
			},
		}
	})
	return singletonBufferPool
}

func newBuffer() interface{} {
	return make([]byte, config.MaxDatagramSize)
}

func (bp *bufferPool) Put(buffer []byte) {
	bp.pool.Put(buffer[:config.MaxDatagramSize])
}

func (bp *bufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}
