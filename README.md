# 日志库
默认是打印控制台等级是debug，可以选择目录进行文件日志输出

# 需要优化问题
1. [x] `func (p *Pool) Put(buf Buffer)`占用过多内存
2. [ ] 系统函数锁占用时间

# 设计
```mermaid
graph LR
    A(WriteChannel) --> B(RecvChannel)
    B --> C(Pack) 
    C --> D(WriteFile)
```