// Package bufpool bufpool

package bufpool

import (
    "reflect"
    "sync"
    "time"
    "unsafe"
)

// LogMsg
type LogMsg struct {
    When time.Time
    Buf  *Buffer
}

// SPool
type SPool struct {
    size int
    p    *sync.Pool
}

// NewSPool new SPool
func NewSPool(size int, pool IPool) *SPool {
    return &SPool{
        p: &sync.Pool{
            New: func() interface{} {
                ret := make([]LogMsg, size)
                for i := range ret {
                    ret[i].Buf = pool.Get()
                }
                return ret
            },
        },
        size: size,
    }
}

func (p *SPool) Get() []LogMsg {
    buf := p.p.Get().([]LogMsg)
    return buf
}

func (p *SPool) Put(buf []LogMsg) {
    (*reflect.SliceHeader)(unsafe.Pointer(&buf)).Len = p.size
    p.p.Put(buf)
}
