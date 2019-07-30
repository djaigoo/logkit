// Package logkit logkit

// +build !windows

package logkit

var (
    // 前景色
    black   = []byte{27, 91, 57, 55, 59, 51, 48, 109}
    red     = []byte{27, 91, 57, 55, 59, 51, 49, 109}
    green   = []byte{27, 91, 57, 55, 59, 51, 50, 109}
    yellow  = []byte{27, 91, 57, 55, 59, 51, 51, 109}
    blue    = []byte{27, 91, 57, 55, 59, 51, 52, 109}
    magenta = []byte{27, 91, 57, 55, 59, 51, 53, 109}
    cyan    = []byte{27, 91, 57, 55, 59, 51, 54, 109}
    white   = []byte{27, 91, 57, 48, 59, 51, 55, 109}
    
    // 背景色
    bblack   = []byte{27, 91, 57, 55, 59, 52, 48, 109}
    bred     = []byte{27, 91, 57, 55, 59, 52, 49, 109}
    bgreen   = []byte{27, 91, 57, 55, 59, 52, 50, 109}
    byellow  = []byte{27, 91, 57, 55, 59, 52, 51, 109}
    bblue    = []byte{27, 91, 57, 55, 59, 52, 52, 109}
    bmagenta = []byte{27, 91, 57, 55, 59, 52, 53, 109}
    bcyan    = []byte{27, 91, 57, 55, 59, 52, 54, 109}
    bwhite   = []byte{27, 91, 57, 48, 59, 52, 55, 109}
    
    reset = []byte{0x1b, 0x5b, 0x30, 0x6d}
)

// brush is a color join function
type brush func(text []byte) []byte

// newBrush return a fix color Brush
func newBrush(color []byte) brush {
    return func(text []byte) []byte {
        return append(color, append(text, reset...)...)
    }
}

var colors = []brush{
    newBrush(white),  // Default
    newBrush(bblue),  // Debug
    newBrush(yellow), // Warning
    newBrush(blue),   // Informational
    newBrush(red),    // Error
    newBrush(green),  // JSON
    
    newBrush(white),   // Emergency
    newBrush(cyan),    // Alert
    newBrush(magenta), // Critical
    newBrush(green),   // Notice
}
