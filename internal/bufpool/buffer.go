// Package bufpool bufpool

package bufpool

// Buffer
type Buffer struct {
    bs []byte
}

// Reset Reset
func (b *Buffer) Reset() {
    b.bs = b.bs[:0]
}

// WriteString WriteString
func (b *Buffer) WriteString(s string) {
    b.bs = append(b.bs, s...)
}

// Write Write
func (b *Buffer) Write(s []byte) {
    b.bs = append(b.bs, s...)
}

// WriteByge WriteByge
func (b *Buffer) WriteByte(c byte) {
    b.bs = append(b.bs, c)
}

// String String
func (b *Buffer) String() string {
    return string(b.bs)
}

// Bytes Bytes
func (b *Buffer) Bytes() []byte {
    return b.bs[:len(b.bs)]
}

// Len Len
func (b *Buffer) Len() int {
    return len(b.bs)
}

// Cap Cap
func (b *Buffer) Cap() int {
    return cap(b.bs)
}
