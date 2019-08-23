// Package logkit logkit

package logkit

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "runtime/debug"
)

// Logkit
type Logkit struct {
    LogWriter
    level Level
}

var lk *Logkit

func init() {
    ConsoleLog(LevelDebug)
}

func ConsoleLog(logLevel Level) {
    lk = new(Logkit)
    lk.level = logLevel
    lk.LogWriter = newConsoleLog(logLevel)
}

func FileLog(basePath, projName string, logLevel Level) *fileObj {
    lk = new(Logkit)
    lk.level = logLevel
    lk.LogWriter = newFileLog(basePath, projName)
    return &fileObj{}
}

func level() Level {
    return lk.level
}

func Exit() {
    tt := &Logkit{
        level:     lk.level,
        LogWriter: newConsoleLog(lk.level),
    }
    lk, tt = tt, lk
    tt.Exit()
}

func Flush() {
    lk.Flush()
}

func Debug(args ...interface{}) {
    if level() <= LevelDebug {
        lk.Write(LevelDebug, []byte(fmt.Sprint(args...)))
    }
}

func Debugf(format string, args ...interface{}) {
    if level() <= LevelDebug {
        lk.Write(LevelDebug, []byte(fmt.Sprintf(format, args...)))
    }
}

func Info(args ...interface{}) {
    if level() <= LevelInfo {
        lk.Write(LevelInfo, []byte(fmt.Sprint(args...)))
    }
}

func Infof(format string, args ...interface{}) {
    if level() <= LevelInfo {
        lk.Write(LevelInfo, []byte(fmt.Sprintf(format, args...)))
    }
}

func Warn(args ...interface{}) {
    if level() <= LevelWarn {
        lk.Write(LevelWarn, []byte(fmt.Sprint(args...)))
    }
}

func Warnf(format string, args ...interface{}) {
    if level() <= LevelWarn {
        lk.Write(LevelWarn, []byte(fmt.Sprintf(format, args...)))
    }
}

func Error(args ...interface{}) {
    if level() <= LevelError {
        lk.Write(LevelError, []byte(fmt.Sprint(args...)))
    }
}

func Errorf(format string, args ...interface{}) {
    if level() <= LevelError {
        lk.Write(LevelError, []byte(fmt.Sprintf(format, args...)))
    }
}

func JSON(v interface{}) error {
    if level() <= LevelJson {
        str, err := json.Marshal(v)
        if err != nil {
            return errors.New("JSON data is empty")
        }
        if len(str) == 0 {
            return nil
        }
        lk.Write(LevelJson, str)
    }
    return nil
}

// Trace console print func trace
func Trace() {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    trace := debug.Stack()
    head := trace[:bytes.Index(trace, []byte{'\n'})+1]
    for i := 0; i < 5; i++ {
        trace = trace[bytes.Index(trace, []byte{'\n'})+1:]
    }
    lk.Write(LevelTrace, append(head, trace...))
}
