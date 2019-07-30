// Package logkit logkit

package logkit

// LogWriter
type LogWriter interface {
    Exit()
    Flush()
    Write(level Level, msg []byte)
}
