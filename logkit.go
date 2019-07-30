// Package logkit logkit

package logkit

import (
    "encoding/json"
    "errors"
    "fmt"
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
    ConsoleLog(lk.level)
    lk.Exit()
}

func Flush() {
    lk.Flush()
}

func Debug(str string) {
    if level() <= LevelDebug {
        lk.Write(LevelDebug, []byte(str))
    }
}

func Debugs(args ...interface{}) {
    if level() <= LevelDebug {
        lk.Write(LevelDebug, []byte(fmt.Sprint(args...)))
    }
}

func Debugf(format string, args ...interface{}) {
    if level() <= LevelDebug {
        lk.Write(LevelDebug, []byte(fmt.Sprintf(format, args...)))
    }
}

func Info(str string) {
    if level() <= LevelInfo {
        lk.Write(LevelInfo, []byte(str))
    }
}

func Infos(args ...interface{}) {
    if level() <= LevelInfo {
        lk.Write(LevelInfo, []byte(fmt.Sprint(args...)))
    }
}

func Infof(format string, args ...interface{}) {
    if level() <= LevelInfo {
        lk.Write(LevelInfo, []byte(fmt.Sprintf(format, args...)))
    }
}

func Warn(str string) {
    if level() <= LevelWarn {
        lk.Write(LevelWarn, []byte(str))
    }
}

func Warns(args ...interface{}) {
    if level() <= LevelWarn {
        lk.Write(LevelWarn, []byte(fmt.Sprint(args...)))
    }
}

func Warnf(format string, args ...interface{}) {
    if level() <= LevelWarn {
        lk.Write(LevelWarn, []byte(fmt.Sprintf(format, args...)))
    }
}

func Error(str string) {
    if level() <= LevelError {
        lk.Write(LevelError, []byte(str))
    }
}

func Errors(args ...interface{}) {
    if level() <= LevelError {
        lk.Write(LevelError, []byte(fmt.Sprint(args...)))
    }
}

func Errorf(format string, args ...interface{}) {
    if level() <= LevelError {
        lk.Write(LevelError, []byte(fmt.Sprintf(format, args...)))
    }
}

// JSON 打印日志格式不带时间戳，每行一个json串
func JSON(v interface{}) error {
    str, err := json.Marshal(v)
    if err != nil {
        return errors.New("JSON data is empty")
    }
    lk.Write(LevelJson, str)
    return nil
}
