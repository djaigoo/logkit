// Package logkit logkit

package logkit

import (
    "encoding/json"
    "io"
    "os"
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
    when := getNow()
    str := pool.Get()
    switch level {
    case LevelJson:
        lt := &logTime{formatTimeHeaderString(when)}
        timestamp, _ := json.Marshal(lt)
        str.Write(timestamp[:len(timestamp)-1])
        str.WriteByte(',')
        str.Write(msg[1:])
    case LevelTrace:
        str.Write(msg)
    default:
        str.Write(formatTimeHeader(when))
        str.WriteByte(' ')
        str.Write(logName[level])
        str.WriteByte(' ')
        str.Write(msg)
    }
    if c, ok := colors[level]; ok {
        cl.writer.Write(append(c(str.Bytes()), '\n'))
    } else {
        cl.writer.Write(append(str.Bytes(), '\n'))
    }
    pool.Put(str)
}
