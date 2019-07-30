// Package logkit logkit

// +build windows

package logkit

var (
    w32Green   = []byte{27, 91, 52, 50, 109}
    w32White   = []byte{27, 91, 52, 55, 109}
    w32Yellow  = []byte{27, 91, 52, 51, 109}
    w32Red     = []byte{27, 91, 52, 49, 109}
    w32Blue    = []byte{27, 91, 52, 52, 109}
    w32Magenta = []byte{27, 91, 52, 53, 109}
    w32Cyan    = []byte{27, 91, 52, 54, 109}
)

var colors = []brush{
    newBrush(),
    newBrush(),
    newBrush(),
    newBrush(),
    newBrush(),
    newBrush(),
    
    newBrush(),
    newBrush(),
    newBrush(),
    newBrush(),
}

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
func newBrush() brush {
    return func(text string) string {
        return text
    }
}
