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

var colors = map[Level]brush{
    LevelDefault:   newBrush(), // Default
    LevelDebug:     newBrush(), // Debug
    LevelWarn:      newBrush(), // Warning
    LevelInfo:      newBrush(), // Informational
    LevelError:     newBrush(), // Error
    levelDivision:  newBrush(), // Default
    LevelJson:      newBrush(), // JSON
    LevelTrace:     newBrush(), // Trace
    LevelEmergency: newBrush(), // Emergency
    LevelAlert:     newBrush(), // Alert
    LevelCritical:  newBrush(), // Critical
    LevelNotice:    newBrush(), // Notice
}

// brush is a color join function
type brush func(string) string

// newBrush return a fix color Brush
func newBrush() brush {
    return func(text string) string {
        return text
    }
}
