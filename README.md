# 日志库
logkit是一个简单的golang日志库。
logkit支持两种日志输出模式：默认是输出控制台日志；还可以将日志输出到文件。
控制台支持输出等级：
```text
debug：调试日志
warn：警告日志
info：信息日志
notice：通知日志
json：json串日志
critical：不建议操作日志
alert：警报日志
error：错误日志
emergency：紧急日志
trace：函数调用路径日志
```
文件支持输出等级：
```text
debug：调试日志
warn：警告日志
info：信息日志
json：json串日志
error：错误日志
```
文件支持两种切割方式：按照时间（每种等级一个日志文件）；按文件大小（所有日志输入一个文件）切割。
其中json日志输出格式仅支持第一种方式输出。

# 需要优化问题
1. [x] `func (p *Pool) Put(buf Buffer)`占用过多内存
2. [ ] 系统函数锁占用时间
