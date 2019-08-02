// Package bufpool bufpool

package bufpool

// IPool
type IPool interface {
    Get() *Buffer
    Put(*Buffer)
}
