# logkit
logkit是一个简单的文件日志库。

优点：
* 调用简单
* 快速写日志
* 轻量级

# 控制台日志
直接通过包名引用即可，默认日志等级是debug，例如
```go
logkit.Debug("hello world")
logkit.Debugf("hello %s", name)
```
如果需要在输出的时候忽略掉debug消息，则可以在使用前设置输出日志等级限制
```go
func init() {
    logkit.ConsoleLog(LevelInfo)
}
```

# 文件日志
如果需要使用文件日志，则必须初始化文件日志对象
```go
func init() {
    logkit.FileLog("base_path", "proj_path", logkit.LevelInfo)
}
```
初始化后使用方式和控制台一样

# 日志等级
在levelDivision之上的都是可以输出的日志等级，下面的都是预留等级。
```go
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
```
在初始化的时候设置日志等级，会忽略设置值上面的，例如
```go
    logkit.ConsoleLog(LevelInfo) // 不会输出debug和warn日志
```

# 控制台日志颜色
windows还不支持颜色
```go
// +build !windows

var colors = map[Level]brush{
    LevelDefault:   newBrush(white),   // Default
    LevelDebug:     newBrush(bblue),   // Debug
    LevelWarn:      newBrush(yellow),  // Warning
    LevelInfo:      newBrush(blue),    // Informational
    LevelError:     newBrush(red),     // Error
    levelDivision:  newBrush(white),   // Default
    LevelJson:      newBrush(black),   // JSON
    LevelTrace:     newBrush(magenta), // Trace
    LevelEmergency: newBrush(white),   // Emergency
    LevelAlert:     newBrush(cyan),    // Alert
    LevelCritical:  newBrush(magenta), // Critical
    LevelNotice:    newBrush(green),   // Notice
}
```