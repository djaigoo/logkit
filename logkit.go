// Package logkit logkit

package logkit

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "runtime/debug"
    "time"
)

// fileObj 临时对象
type fileObj struct{}

// RotateSize 设置切换文件大小
func (fo *fileObj) SetRotateSize(size uint64) *fileObj {
    if size == 0 {
        return fo
    }
    maxSize = size
    return fo
}

// FlushTime 设置刷新时间
func (fo *fileObj) SetFlushTime(ft time.Duration) *fileObj {
    if ft == 0 {
        return fo
    }
    flushTime = ft
    df, ok := lk.LogWriter.(interface{ changeTicker() })
    if ok {
        df.changeTicker()
    }
    return fo
}

// Logkit
type Logkit struct {
    LogWriter
    level Level
}

var lk *Logkit

func init() {
    ConsoleLog(LevelDefault)
}

func ConsoleLog(logLevel Level) {
    lk = new(Logkit)
    lk.level = logLevel
    lk.LogWriter = newConsoleLog(logLevel)
}

func SingleFileLog(basePath, projName string, logLevel Level) *fileObj {
    lk = new(Logkit)
    lk.level = logLevel
    lk.LogWriter = newSingleFileLog(basePath, projName)
    return &fileObj{}
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

func SetLevel(l Level) {
    lk.level = l
}

func SetLevelString(l string) *fileObj {
    switch l {
    case "debug":
        lk.level = LevelDebug
    case "warn":
        lk.level = LevelWarn
    case "info":
        lk.level = LevelInfo
    case "error":
        lk.level = LevelError
    default:
    
    }
    return &fileObj{}
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

func write(l Level, args ...interface{}) {
    if level() <= l {
        lk.Write(l, []byte(fmt.Sprint(args...)))
    }
}

func writef(l Level, format string, args ...interface{}) {
    if level() <= l {
        lk.Write(l, []byte(fmt.Sprintf(format, args...)))
    }
}

func Debug(args ...interface{}) {
    write(LevelDebug, args...)
}

func Debugf(format string, args ...interface{}) {
    writef(LevelDebug, format, args...)
}

func Info(args ...interface{}) {
    write(LevelInfo, args...)
}

func Infof(format string, args ...interface{}) {
    writef(LevelInfo, format, args...)
}

func Warn(args ...interface{}) {
    write(LevelWarn, args...)
}

func Warnf(format string, args ...interface{}) {
    writef(LevelWarn, format, args...)
}

func Error(args ...interface{}) {
    write(LevelError, args...)
}

func Errorf(format string, args ...interface{}) {
    writef(LevelError, format, args...)
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

func Emergency(args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    write(LevelEmergency, args...)
}

func Emergencyf(format string, args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    writef(LevelEmergency, format, args...)
}

func Alert(args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    write(LevelAlert, args...)
}

func Alertf(format string, args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    writef(LevelAlert, format, args...)
}

func Critical(args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    write(LevelCritical, args...)
}

func Criticalf(format string, args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    writef(LevelCritical, format, args...)
}

func Notice(args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    write(LevelNotice, args...)
    
}

func Noticef(format string, args ...interface{}) {
    if _, ok := lk.LogWriter.(*consoleLog); !ok {
        return
    }
    writef(LevelNotice, format, args...)
}
