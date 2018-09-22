package model

import (
	"sync"

	"github.com/bonnetn/srcds_proxy/proxy/config"
)

// BufferPool is a wrapper around sync.Pool that is type-safe.
type BufferPool struct {
	pool sync.Pool
}

var (
	singletonBufferPool *BufferPool
	once                sync.Once
)

// GetBufferPool returns the BufferPool (singleton).
func GetBufferPool() *BufferPool {
	once.Do(func() {
		singletonBufferPool = &BufferPool{
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

// Put adds a new buffer in the pool to be re-used.
func (bp *BufferPool) Put(buffer []byte) {
	bp.pool.Put(buffer[:config.MaxDatagramSize])
}

// Get removes a buffer from the pool.
func (bp *BufferPool) Get() []byte {
	return bp.pool.Get().([]byte)
}
