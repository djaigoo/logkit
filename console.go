// Package logkit logkit

package logkit

import (
    "bytes"
    "io"
    "os"
    "time"
)

// consoleLog
type consoleLog struct {
    level  Level
    writer io.Writer
}

// newConsoleLog new ConsoleLog
func newConsoleLog(level Level) *consoleLog {
    return &consoleLog{
        level:  level,
        writer: os.Stdout,
    }
}

func (cl *consoleLog) Exit() {
}

func (cl *consoleLog) Flush() {
}

func (cl *consoleLog) Write(level Level, msg []byte) {
    if level == LevelJson {
        cl.writer.Write(append(colors[level](msg), '\n'))
        return
    }
    when := time.Now()
    str := bytes.NewBuffer(formatTimeHeader(when))
    str.WriteByte(' ')
    str.Write(logName[level])
    str.WriteByte(' ')
    str.Write(msg)
    
    cl.writer.Write(append(colors[level](str.Bytes()), '\n'))
}
