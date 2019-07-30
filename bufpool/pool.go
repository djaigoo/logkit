// Package bufpool bufpool

package bufpool

import (
    "sync"
)

// Pool
type Pool struct {
    p *sync.Pool
}

// NewPool constructs a new Pool.
func NewPool(size int) *Pool {
    return &Pool{p: &sync.Pool{
        New: func() interface{} {
            return &Buffer{bs: make([]byte, 0, size)}
        },
    }}
}

// Get retrieves a Buffer from the Pool, creating one if necessary.
func (p *Pool) Get() *Buffer {
    buf := p.p.Get().(*Buffer)
    buf.Reset()
    return buf
}

func (p *Pool) Put(buf *Buffer) {
    p.p.Put(buf)
    
}
