// Package logkit logkit

package logkit

import (
    "bufio"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "time"
    
    "gitlab.66xue.com/daihao/logkit/bufpool"
)

const (
    levelCount    = int(levelDivision) // log level count
    bufLogSize    = 2048               // buf LogMsg size
    preLogLen     = 64                 // 预计日志单条长度
    bufLogHeadLen = 28                 // 每条日志的基本长度
)

var (
    zeroTime, _ = time.Parse("2006-01-02 15:04:05", "01-01-01 00:00:00")
)

var (
    maxSize     uint64 = 1024 * 1024 * 1024 * 5       // single file size limit 50MB
    bufferSize         = 256 * 1024              // bufio buf size
    flushTime          = 1500 * time.Millisecond // flush write file
    fileLogSize        = 256                     // fileLog channel size
    bufChanSize        = 512                     // bufferWriter chan LogMsg size
)

var (
    pool  = bufpool.NewPool(preLogLen)
    spool = bufpool.NewSPool(bufLogSize, pool)
)

// fileObj 临时对象
type fileObj struct{}

// RotateSize 设置切换文件大小
func (fo *fileObj) SetRotateSize(size uint64) *fileObj {
    maxSize = size
    return fo
}

// FlushTime 设置刷新时间
func (fo *fileObj) SetFlushTime(ft time.Duration) *fileObj {
    flushTime = ft
    return fo
}

// fileLog
type fileLog struct {
    mch    [levelCount]chan []byte   // chan LogMsg
    bw     [levelCount]*bufferWriter // write file object
    ticker *time.Ticker              // ticker flush
    wg     *sync.WaitGroup           // wait all buffer writer work over
    handle []func(level Level, msg []byte)
}

// newFileLog new FileLog
func newFileLog(basePath, logName string) *fileLog {
    mch := [levelCount]chan []byte{}
    for i := 0; i < levelCount; i++ {
        mch[i] = make(chan []byte, fileLogSize)
    }
    bw := [levelCount]*bufferWriter{}
    for i := 0; i < levelCount; i++ {
        bw[i] = newBufferWriter(basePath, logName, Level(i))
    }
    ret := &fileLog{
        mch:    mch,
        bw:     bw,
        ticker: time.NewTicker(flushTime),
        wg:     new(sync.WaitGroup),
    }
    
    ret.transLog()
    go ret.flushTicker()
    
    return ret
}

func (fl *fileLog) Exit() {
    fl.ticker.Stop()
    
    for i := range fl.mch {
        close(fl.mch[i])
    }
    fl.wg.Wait()
    for i := range fl.bw {
        fl.bw[i].Close()
    }
}

func (fl *fileLog) Flush() {
    lm := bufpool.LogMsg{When: zeroTime, Buf: pool.Get()}
    for _, bw := range fl.bw {
        bw.WriteLog(lm)
    }
}

func (fl *fileLog) Write(level Level, msg []byte) {
    fl.mch[level] <- msg
}

func (fl *fileLog) flushTicker() {
    lm := bufpool.LogMsg{When: zeroTime, Buf: pool.Get()}
    for range fl.ticker.C {
        for i := 0; i < levelCount; i++ {
            fl.bw[i].WriteLog(lm)
        }
    }
}

func (fl *fileLog) transLog() {
    for i := range fl.mch {
        i := i
        fl.wg.Add(1)
        go func() {
            var msg []byte
            for msg = range fl.mch[i] {
                // get pool
                buf := pool.Get()
                buf.Write(msg)
                buf.WriteByte('\n')
                lm := bufpool.LogMsg{When: getNow(), Buf: buf}
                fl.bw[i].WriteLog(lm)
            }
            fl.wg.Done()
        }()
    }
}

// bufferWriter
type bufferWriter struct {
    *bufio.Writer
    buf      []bufpool.LogMsg      // buffer
    offset   int                   // buffer offset
    lmch     chan []bufpool.LogMsg // LogMsg chan
    basePath string                // 基础路径
    logName  string                // 项目名称
    file     *os.File              // 日志文件指针
    level    Level                 // 日志等级
    slot     int                   // 文件分隔计数
    nbytes   uint64                // 统计写入字符数
    sm       *sync.Mutex           // 修改buf加锁
    wg       *sync.WaitGroup       // 等待日志写完
    lastHour int
}

// newBufferWriter new bufferWriter
func newBufferWriter(base, log string, level Level) *bufferWriter {
    ret := &bufferWriter{
        buf:      spool.Get(), // make([]bufpool.LogMsg, bufLogSize),
        lmch:     make(chan []bufpool.LogMsg, bufChanSize),
        basePath: base,
        logName:  log,
        level:    level,
        sm:       new(sync.Mutex),
        wg:       new(sync.WaitGroup),
    }
    // 续接上次日志文件
    tn := time.Now()
    year, month, day := tn.Date()
    logDir := filepath.Join(base, log, fmt.Sprintf("%04d%02d", year, month), fmt.Sprintf("%2d", day))
    info, _ := os.Stat(logDir)
    if info != nil && info.IsDir() {
        prefix := fmt.Sprintf("%s-%02d-", levelName[level], tn.Hour())
        files, _ := ioutil.ReadDir(logDir)
        name := ""
        for _, f := range files {
            if strings.HasPrefix(f.Name(), prefix) {
                if len(name) <= len(f.Name()) {
                    name = f.Name()
                }
            }
        }
        if name != "" {
            tpath := filepath.Join(logDir, name)
            info, _ = os.Stat(tpath)
            ret.nbytes = uint64(info.Size())
            tags := strings.Split(name, "-")
            if len(tags) == 3 {
                tags = strings.Split(tags[2], ".")
                ret.slot, _ = strconv.Atoi(tags[0])
            } else {
                ret.slot = 0
            }
            f, _ := os.OpenFile(tpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
            ret.file = f
            ret.Writer = bufio.NewWriterSize(ret.file, bufferSize)
        }
    }
    ret.wg.Add(1)
    go func() {
        ret.writeFile()
        ret.wg.Done()
    }()
    return ret
}

// Close Close
func (bw *bufferWriter) Close() {
    // flush
    bw.WriteLog(bufpool.LogMsg{When: zeroTime, Buf: pool.Get()})
    close(bw.lmch)
    bw.wg.Wait()
    bw.file.Close()
}

// WriteLog WriteLog
func (bw *bufferWriter) WriteLog(lm bufpool.LogMsg) () {
    // chan zero is sign to Flush
    if lm.When.IsZero() {
        bw.Flush()
        return
    }
    offset := 0
    
    bw.sm.Lock()
    pool.Put(bw.buf[bw.offset].Buf)
    bw.buf[bw.offset] = lm
    bw.offset++
    offset = bw.offset
    bw.sm.Unlock()
    
    if offset == bufLogSize {
        bw.Flush()
    }
}

// Flush buffer to file
func (bw *bufferWriter) Flush() {
    if bw.offset == 0 {
        return
    }
    buf := spool.Get()
    bw.sm.Lock()
    bw.buf, buf = buf, bw.buf
    offset := bw.offset
    bw.offset = 0
    bw.sm.Unlock()
    bw.lmch <- buf[:offset]
}

func (bw *bufferWriter) writeFile() {
    var datas []bufpool.LogMsg
    for datas = range bw.lmch {
        if len(datas) == 0 {
            continue
        }
        if bw.file == nil || bw.Writer == nil {
            tn := time.Now()
            bw.rotateFile(0, tn)
        }
        
        i := 0
        if bw.lastHour != datas[0].When.Hour() {
            bw.lastHour = datas[0].When.Hour()
            bw.rotateFile(0, datas[0].When)
            for i = range datas {
                bw.writeData(datas[i])
                bw.nbytes += uint64(datas[i].Buf.Len() + bufLogHeadLen)
            }
        } else if bw.lastHour == datas[len(datas)-1].When.Hour() {
            // 粗略分割
            if bw.nbytes >= maxSize {
                bw.rotateFile(bw.slot+1, datas[len(datas)-1].When)
            }
            for i = range datas {
                bw.writeData(datas[i])
                bw.nbytes += uint64(datas[i].Buf.Len() + bufLogHeadLen)
            }
        } else {
            once := sync.Once{}
            for _, data := range datas {
                if data.When.Hour() != bw.lastHour {
                    bw.writeData(data)
                    continue
                }
                once.Do(func() {
                    bw.lastHour = data.When.Hour()
                    bw.rotateFile(0, data.When)
                    bw.nbytes = 0
                })
                bw.writeData(data)
                bw.nbytes += uint64(data.Buf.Len() + bufLogHeadLen)
            }
        }
        spool.Put(datas)
        bw.Writer.Flush()
    }
}

func (bw *bufferWriter) writeData(data bufpool.LogMsg) {
    bw.Writer.WriteByte('[')
    bw.Writer.Write(formatTimeHeader(data.When))
    bw.Writer.WriteByte(']')
    bw.Writer.WriteByte(' ')
    bw.Writer.Write(data.Buf.Bytes())
}

func (bw *bufferWriter) rotateFile(slot int, when time.Time) {
    if bw.file != nil {
        tWriter := bw.Writer
        tFile := bw.file
        go func() {
            tWriter.Flush()
            tFile.Close()
        }()
    }
    file, err := createFile(bw.basePath, bw.logName, bw.level, slot, when)
    if err != nil {
        panic(err.Error())
    }
    bw.file = file
    bw.nbytes = 0
    bw.slot = slot
    
    bw.Writer = bufio.NewWriterSize(bw.file, bufferSize)
}

// Sync sync disk so slow
func (bw *bufferWriter) Sync() error {
    return bw.file.Sync()
}

func createFile(basePath, logName string, level Level, slot int, when time.Time) (*os.File, error) {
    year, month, day := when.Date()
    if !strings.HasSuffix(basePath, "/") {
        basePath += "/"
    }
    logDir := filepath.Join(basePath, logName, fmt.Sprintf("%04d%02d", year, month), fmt.Sprintf("%02d", day))
    err := os.MkdirAll(logDir, os.ModePerm)
    if err != nil {
        return nil, fmt.Errorf("logkit: cannot create log: %v", err)
    }
    var logFile string
    if slot <= 0 {
        logFile = fmt.Sprintf("%s-%02d.log", levelName[level], when.Hour())
    } else {
        logFile = fmt.Sprintf("%s-%02d-%d.log", levelName[level], when.Hour(), slot)
    }
    fname := filepath.Join(logDir, logFile)
    f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return nil, fmt.Errorf("logkit: cannot open log file: %v", err)
    }
    return f, nil
}
