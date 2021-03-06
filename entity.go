// Package logkit logkit

package logkit

import (
    "github.com/djaigoo/logkit/internal/bufpool"
)

type Level uint8

const (
    LevelDefault = Level(iota)
    LevelDebug
    LevelWarn
    LevelInfo
    LevelError
    LevelJson
    levelDivision
    
    LevelTrace
    LevelEmergency
    LevelAlert
    LevelCritical
    LevelNotice
    levelMax
    LevelNon = levelMax // not use log
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
        
        LevelEmergency: {'[', 'M', ']'},
        LevelAlert:     {'[', 'A', ']'},
        LevelCritical:  {'[', 'C', ']'},
        LevelNotice:    {'[', 'N', ']'},
    }
)

var (
    pool  = bufpool.NewPool(preLogLen)
    spool = bufpool.NewSPool(bufLogSize, pool)
)
