package logkit

import (
    "bufio"
    "fmt"
    "gitlab.66xue.com/daihao/logkit/internal/bufpool"
    "os"
    "path/filepath"
    "runtime"
    "sync"
    "time"
)

const (
    logFileName = "2006-01-02_15:04:05.log"
)

// singleFile
type singleFile struct {
    base   string
    name   string
    ticker *time.Ticker
    wg     *sync.WaitGroup
    mch    chan *bufpool.Buffer // chan LogMsg
    bw     *sbufferWriter       // write file object
}

func newSingleFileLog(basePath, logName string) *singleFile {
    mch := make(chan *bufpool.Buffer, fileLogSize)
    if basePath == "" {
        basePath = "."
    }
    bw := newSBufferWriter(basePath, logName)
    ret := &singleFile{
        mch:    mch,
        bw:     bw,
        ticker: time.NewTicker(flushTime),
        wg:     new(sync.WaitGroup),
    }
    
    ret.transLog()
    go ret.flushTicker()
    
    return ret
}

// changeTicker changeTicker
func (sf *singleFile) changeTicker() {
    sf.ticker.Stop()
    sf.ticker = time.NewTicker(flushTime)
    go sf.flushTicker()
}

func (sf *singleFile) Exit() {
    sf.ticker.Stop()
    close(sf.mch)
    sf.wg.Wait()
    sf.bw.Close()
}

func (sf *singleFile) Flush() {
    lm := bufpool.LogMsg{When: zeroTime, Buf: pool.Get()}
    sf.bw.WriteLog(lm)
}

func (sf *singleFile) Write(level Level, msg []byte) {
    buf := pool.Get()
    buf.Write(logName[level])
    buf.WriteByte(' ')
    buf.Write(msg)
    buf.WriteByte('\n')
    sf.mch <- buf
}

func (sf *singleFile) flushTicker() {
    lm := bufpool.LogMsg{When: zeroTime, Buf: pool.Get()}
    for range sf.ticker.C {
        sf.bw.WriteLog(lm)
    }
}

func (sf *singleFile) transLog() {
    sf.wg.Add(1)
    go func() {
        var msg *bufpool.Buffer
        for msg = range sf.mch {
            lm := bufpool.LogMsg{When: getNow(), Buf: msg}
            sf.bw.WriteLog(lm)
        }
        sf.wg.Done()
    }()
}

type sbufferWriter struct {
    *bufio.Writer
    buf      []bufpool.LogMsg      // buffer
    offset   int                   // buffer offset
    lmch     chan []bufpool.LogMsg // LogMsg chan
    basePath string                // 基础路径
    logName  string                // 项目名称
    file     *os.File              // 日志文件指针
    nbytes   uint64                // 统计写入字符数
    sm       *sync.Mutex           // 修改buf加锁
    wg       *sync.WaitGroup
}

// newBufferWriter new bufferWriter
func newSBufferWriter(base, log string) *sbufferWriter {
    ret := &sbufferWriter{
        buf:      spool.Get(), // make([]bufpool.LogMsg, bufLogSize),
        lmch:     make(chan []bufpool.LogMsg, bufChanSize),
        basePath: base,
        logName:  log,
        sm:       new(sync.Mutex),
        wg:       new(sync.WaitGroup),
    }
    var err error
    ret.file, err = createFileSingle(base, log, getNow())
    if err != nil {
        panic(err.Error())
        return ret
    }
    ret.Writer = bufio.NewWriterSize(ret.file, bufferSize)
    ret.wg.Add(1)
    go ret.writeFile(bufLogHeadLen, ret.writeData)
    return ret
}

// Close Close
func (bw *sbufferWriter) Close() {
    // flush
    bw.WriteLog(bufpool.LogMsg{When: zeroTime, Buf: pool.Get()})
    close(bw.lmch)
    bw.wg.Wait()
    bw.file.Close()
}

// WriteLog WriteLog
func (bw *sbufferWriter) WriteLog(lm bufpool.LogMsg) () {
    // chan zero is sign to Flush
    if lm.When.IsZero() {
        bw.Flush()
        return
    }
    offset := 0
    
    bw.sm.Lock()
    bw.buf[bw.offset] = lm
    bw.offset++
    offset = bw.offset
    bw.sm.Unlock()
    
    if offset == bufLogSize {
        bw.Flush()
    }
}

// Flush buffer to file
func (bw *sbufferWriter) Flush() {
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

func (bw *sbufferWriter) writeFile(slen int, writeData func(data bufpool.LogMsg)) {
    defer bw.wg.Done()
    runtime.LockOSThread()
    var datas []bufpool.LogMsg
    for datas = range bw.lmch {
        if len(datas) == 0 {
            continue
        }
        if bw.file == nil || bw.Writer == nil {
            tn := getNow()
            bw.rotateFile(tn)
        }
        
        if bw.nbytes >= maxSize {
            bw.rotateFile(getNow())
        }
        for i := range datas {
            writeData(datas[i])
            bw.nbytes += uint64(datas[i].Buf.Len() + slen)
        }
        spool.Put(datas)
        bw.Writer.Flush()
    }
    
    runtime.UnlockOSThread()
}

func (bw *sbufferWriter) writeData(data bufpool.LogMsg) {
    buf := pool.Get()
    buf.WriteByte('[')
    buf.Write(formatTime(data.When))
    buf.WriteByte(']')
    buf.WriteByte(' ')
    bw.Writer.Write(buf.Bytes())
    pool.Put(buf)
    bw.Writer.Write(data.Buf.Bytes())
    pool.Put(data.Buf)
}

func (bw *sbufferWriter) rotateFile(when time.Time) {
    if bw.file != nil {
        tWriter := bw.Writer
        tFile := bw.file
        go func() {
            tWriter.Flush()
            tFile.Close()
        }()
    }
    file, err := createFileSingle(bw.basePath, bw.logName, when)
    if err != nil {
        // panic(err.Error())
        return
    }
    bw.file = file
    bw.nbytes = 0
    
    bw.Writer = bufio.NewWriterSize(bw.file, bufferSize)
}

// Sync sync disk so slow
func (bw *sbufferWriter) Sync() error {
    return bw.file.Sync()
}

func createFileSingle(basePath, logName string, when time.Time) (*os.File, error) {
    filename := when.Format(logFileName)
    lp := filepath.Join(basePath, logName)
    err := os.MkdirAll(lp, os.ModePerm)
    if err != nil {
        return nil, fmt.Errorf("logkit: cannot create log: %v", err)
    }
    lp = filepath.Join(lp, filename)
    f, err := os.OpenFile(lp, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return nil, fmt.Errorf("logkit: cannot open log file: %v", err)
    }
    return f, nil
}
