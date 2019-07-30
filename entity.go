// Package logkit logkit

package logkit

type Level uint8

const (
    LevelDefault = Level(iota)
    LevelDebug
    LevelWarn
    LevelInfo
    LevelError
    levelDivision
    LevelJson
    levelMax
)

var (
    // log file name map
    levelName = map[Level]string{
        LevelDefault: "default",
        LevelDebug:   "debug",
        LevelWarn:    "warn",
        LevelInfo:    "info",
        LevelError:   "error",
        LevelJson:    "json",
    }
    // log file row tag
    logName = map[Level][]byte{
        LevelDefault: {},
        LevelDebug:   {'[', 'D', ']'},
        LevelWarn:    {'[', 'W', ']'},
        LevelInfo:    {'[', 'I', ']'},
        LevelError:   {'[', 'E', ']'},
    }
)
