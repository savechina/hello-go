# 高级并发 (Advanced Concurrency)

## 开篇故事

想象你在一个共享厨房做饭。如果只有一个人，随便用哪个锅都行。但如果有 100 个人同时要做饭，问题就来了：

- 如果两个人同时往同一个锅里加菜 → 菜洒一地（数据混乱）
- 如果有人只看菜谱（不碰锅），其实可以多人同时看
- 如果只是数一下有多少个盘子，不需要抢锅，用个计数器就行

在 Go 程序中，goroutine 就像这些厨师。`sync.Mutex`、`sync.RWMutex` 和 `sync/atomic` 就是管理共享厨房的规则。选错工具，程序就会"数据竞争"（data race）——这是最难调试的 bug 之一。

## 本章适合谁

- ✅ 已经用过 goroutine，但遇到"goroutine 改了数据，另一个 goroutine 读不到"的问题
- ✅ 用过 channel，但发现某些场景用锁更方便（如保护缓存、计数器）
- ✅ 想理解 Mutex、RWMutex、atomic 的区别，知道何时用哪个
- ✅ 遇到过"程序偶尔输出错误结果，但不知道何时发生"的竞态条件

如果你曾经写过 `counter++` 在多个 goroutine 中，然后发现结果"有时候对，有时候不对"，本章必读。

## 你会学到什么

完成本章后，你将能够：

1. **正确使用 Mutex**：用 `Lock()`/`Unlock()` 保护临界区，避免数据竞争
2. **区分 Mutex 和 RWMutex**：理解"读多写少"场景，用 RWMutex 提升性能
3. **掌握 atomic 操作**：用 `sync/atomic` 实现高效计数器
4. **识别竞态条件**：通过代码审查发现缺少保护的共享变量
5. **使用 go test -race**：用竞态检测器验证并发代码的正确性

## 前置要求

在开始之前，请确保你已掌握：

- Go 基础语法（变量、指针、结构体）
- goroutine 的启动和 WaitGroup 的使用
- 理解什么是"共享变量"和"并发修改"
- 了解 channel 的基本用法（可选，但有助于对比）

如果不确定什么是"竞态条件"，可以先阅读《并发基础》章节。

## 第一个例子

让我们从一个最简单的场景开始：100 个 goroutine 同时给一个计数器加 1。

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var mu sync.Mutex
    counter := 0
    var wg sync.WaitGroup
    
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    
    wg.Wait()
    fmt.Printf("计数 = %d (预期 100)\n", counter)
}
```

**运行结果**：
```
计数 = 100 (预期 100)
```

**如果去掉 Mutex**：
```go
// ❌ 没有锁保护
go func() {
    defer wg.Done()
    counter++ // 100 个 goroutine 同时执行这行
}()
```

**可能输出**：`计数 = 87`（每次运行结果不同）

**为什么**：`counter++` 不是原子操作，它分为三步：读取、加 1、写回。当多个 goroutine 同时执行时，会互相覆盖。

## 原理解析

### 1. 什么是竞态条件 (Race Condition)？

**定义**：当两个或多个 goroutine 同时访问同一个变量，且至少有一个是写操作时，就会发生竞态条件。

**直观理解**：想象两个人同时修改同一份文档。A 复制了内容，B 也复制了内容。A 改完保存，B 改完保存。B 的保存会覆盖 A 的修改——A 的工作白做了。

**Go 的例子**：
```go
counter++ // 编译器把它变成三条指令
// 1. MOV counter → register (读取)
// 2. ADD 1 → register (加 1)
// 3. MOV register → counter (写回)
```

当两个 goroutine 同时执行这三条指令时：

```
时间    Goroutine A          Goroutine B
----    ------------         ------------
t1      读取 counter (=0)
t2                           读取 counter (=0)
t3      加 1 (=1)
t4                           加 1 (=1)
t5      写回 counter (=1)
t6                           写回 counter (=1)
```

结果：两次加 1，但 counter 只增加了 1。这就是竞态。

### 2. Mutex 如何解决问题？

`sync.Mutex` 提供了一个"互斥锁"：同一时刻只能有一个 goroutine 持有它。

```go
mu.Lock()   // 尝试加锁，如果已被占用则等待
counter++   // 临界区：只有我能执行
mu.Unlock() // 释放锁，让其他人可以进来
```

**工作流程**：
1. A 调用 `Lock()`，获得锁，进入临界区
2. B 调用 `Lock()`，发现锁被占用，**阻塞等待**
3. A 执行完，调用 `Unlock()`
4. B 被唤醒，获得锁，进入临界区

**关键点**：Mutex 保证了临界区的"互斥访问"，就像卫生间的"使用中"标志。

### 3. RWMutex：读多写少的优化

**问题**：Mutex 太严格了。如果 10 个人都要读文档（不改），为什么要排队？

**解决**：`sync.RWMutex` 区分读锁和写锁：
- `RLock()` / `RUnlock()`：读锁，允许多个读者同时持有
- `Lock()` / `Unlock()`：写锁，独占，与其他所有锁互斥

**使用场景**：
```go
var cache map[string]string
var rwmu sync.RWMutex

// 读操作：可以并发
func get(key string) string {
    rwmu.RLock()
    defer rwmu.RUnlock()
    return cache[key]
}

// 写操作：独占
func set(key, value string) {
    rwmu.Lock()
    defer rwmu.Unlock()
    cache[key] = value
}
```

**性能对比**：
- 100 个 goroutine 只读：RWMutex 比 Mutex 快约 10 倍（因为没有串行化）
- 100 个 goroutine 全写：RWMutex 和 Mutex 性能相当

### 4. atomic：CPU 级别的原子操作

**原理**：`sync/atomic` 使用 CPU 的特殊指令（如 `LOCK XADD`）保证操作原子性，无需操作系统介入。

**对比 Mutex**：
| 特性 | Mutex | atomic |
|------|-------|--------|
| 粒度 | 保护代码块 | 保护单个变量 |
| 性能 | 较慢（涉及系统调用） | 极快（CPU 指令） |
| 适用场景 | 复杂临界区 | 简单计数器 |

**适用类型**：
```go
var (
    i32 int32
    i64 int64
    u32 uint32
    u64 uint64
    ptr unsafe.Pointer
)

atomic.AddInt64(&i64, 1)      // 原子加法
val := atomic.LoadInt64(&i64) // 原子读取
atomic.StoreInt64(&i64, 42)   // 原子写入
```

**⚠️ 限制**：只能用于基本类型，不能保护复杂逻辑。

### 5. WaitGroup：等待多个 goroutine

虽然 WaitGroup 在前面的章节学过，但在这里它是关键配角：

```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done() // 确保 Done() 被调用
        // 做一些事
    }()
}
wg.Wait() // 阻塞，直到所有 goroutine 完成
```

**关键点**：
- `Add(1)` 必须在启动 goroutine **之前**调用
- `Done()` 必须在 goroutine **结束时**调用（用 `defer` 最安全）
- `Wait()` 会阻塞当前 goroutine

## 常见错误

### 错误 1：忘记解锁

```go
// ❌ 错误代码
mu.Lock()
counter++
// 忘了 Unlock()，程序死锁

// 编译器不会报错，但程序会卡住
```

**如何修复**：
```go
// ✅ 修复：用 defer 确保解锁
mu.Lock()
defer mu.Unlock()
counter++
```

**为什么用 defer**：即使临界区内发生 panic，`Unlock()` 也会被调用，避免死锁。

### 错误 2：RWMutex 写操作误用读锁

```go
// ❌ 错误代码
rwmu.RLock()
data["key"] = "value" // 写操作！
rwmu.RUnlock()

// 可能 panic: concurrent map writes
```

**原因**：多个 goroutine 同时持有读锁，同时写 map 会导致 panic。

**修复**：
```go
// ✅ 写操作用写锁
rwmu.Lock()
data["key"] = "value"
rwmu.Unlock()
```

### 错误 3：atomic 类型不匹配

```go
// ❌ 错误代码
var counter int // 普通 int
atomic.AddInt64(&counter, 1) // 编译错误：类型不匹配

// ✅ 修复：用 int64
var counter int64
atomic.AddInt64(&counter, 1)
```

**注意**：`atomic` 包有严格的类型要求，`int` 和 `int64` 不能混用。

## 动手练习

### 练习 1：预测输出

阅读以下代码，预测输出（先自己想，再看答案）：

```go
var counter int
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        counter++
    }()
}
wg.Wait()
fmt.Println(counter)
```

<details>
<summary>点击查看答案</summary>

**答案**：输出不确定，可能是 7、8、9、10 等任意值（通常小于 10）。

**原因**：竞态条件。10 个 goroutine 同时执行 `counter++`，有些操作被覆盖了。

**修复**：加一个 `sync.Mutex` 或使用 `atomic.AddInt64`。
</details>

### 练习 2：修复 RWMutex 误用

以下代码有什么隐患？如何修复？

```go
var cache = make(map[string]int)
var rwmu sync.RWMutex

func update(key string, value int) {
    rwmu.RLock()
    cache[key] = value // 写操作
    rwmu.RUnlock()
}

func get(key string) int {
    rwmu.Lock() // 读操作用写锁
    val := cache[key]
    rwmu.Unlock()
    return val
}
```

<details>
<summary>点击查看答案</summary>

**问题**：
1. `update` 用读锁做写操作 → panic
2. `get` 用写锁做读操作 → 性能浪费（不能并发）

**修复**：
```go
func update(key string, value int) {
    rwmu.Lock()
    defer rwmu.Unlock()
    cache[key] = value
}

func get(key string) int {
    rwmu.RLock()
    defer rwmu.RUnlock()
    return cache[key]
}
```
</details>

### 练习 3：实现线程安全的计数器

用三种方式实现一个计数器（支持并发 Inc() 和 Value()）：
1. 使用 `sync.Mutex`
2. 使用 `sync/atomic`
3. 使用 channel

<details>
<summary>点击查看答案</summary>

```go
// 方式 1: Mutex
type MutexCounter struct {
    mu    sync.Mutex
    value int64
}
func (c *MutexCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
func (c *MutexCounter) Value() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// 方式 2: atomic
type AtomicCounter struct {
    value int64
}
func (c *AtomicCounter) Inc() {
    atomic.AddInt64(&c.value, 1)
}
func (c *AtomicCounter) Value() int64 {
    return atomic.LoadInt64(&c.value)
}

// 方式 3: channel
type ChannelCounter struct {
    ch chan int
}
func NewChannelCounter() *ChannelCounter {
    c := &ChannelCounter{ch: make(chan int, 100)}
    go func() {
        var val int64
        for range c.ch {
            val++
        }
    }()
    return c
}
// 注意：channel 方案读取值较复杂，不适合此场景
```

**推荐**：简单计数器用 atomic，复杂逻辑用 Mutex。
</details>

## 故障排查 (FAQ)

### Q1: 如何检测竞态条件？

**工具**：`go test -race` 或 `go run -race`

**示例**：
```bash
$ go run -race main.go
WARNING: DATA RACE
Read at 0x00c0000140a0 by goroutine 7:
  main.main.func1()
      main.go:15

Previous write at 0x00c0000140a0 by goroutine 6:
  main.main.func1()
      main.go:14
```

**输出解读**：
- 哪些 goroutine 参与了竞争
- 读写操作发生在哪一行代码
- 涉及的内存地址

**建议**：CI 流程中 always 加 `-race` 标志。

### Q2: Mutex 和 channel 应该如何选择？

**原则**：
- **共享状态**（缓存、配置）→ Mutex
- **所有权转移**（任务队列、消息）→ channel
- **简单计数** → atomic

**例子**：
```go
// ✅ Mutex：保护共享缓存
var cache map[string]string
var mu sync.Mutex
func get(key string) { /* 读缓存 */ }

// ✅ channel: 任务分发
jobs := make(chan Job)
go func() {
    for job := range jobs {
        process(job)
    }
}()
```

**Go 的哲学**："不要通过共享内存来通信，而要通过通信来共享内存"。但这条规则有例外——保护共享状态时，Mutex 更直观。

### Q3: RWMutex 的锁升级问题

**问题**：持有读锁时，能升级为写锁吗？

**答案**：**不能**，会导致死锁。

```go
// ❌ 错误代码
rwmu.RLock()
// 发现需要修改...
rwmu.Lock() // 死锁！因为已经有读锁（包括自己的）
```

**正确做法**：先释放读锁，再获取写锁（但要小心在此期间数据被其他人修改）。

## 知识扩展 (选学)

### Mutex 的内部实现

Go 的 `sync.Mutex` 有两种模式：
1. **正常模式**：按 FIFO 顺序唤醒等待者
2. **饥饿模式**：直接 handing off 给等待最久的人，避免"饿死"

当某个 goroutine 等待超过 1ms 时，Mutex 切换到饥饿模式。这是 Go runtime 的自动优化。

### Cond：条件变量

`sync.Cond` 用于更复杂的同步场景（如生产者-消费者）：

```go
cond := sync.NewCond(&sync.Mutex{})
var queue []Item

// 消费者
cond.L.Lock()
for len(queue) == 0 {
    cond.Wait() // 等待，释放锁
}
item := queue[0]
queue = queue[1:]
cond.L.Unlock()

// 生产者
cond.L.Lock()
queue = append(queue, item)
cond.Signal() // 唤醒一个消费者
cond.L.Unlock()
```

### sync.Map：并发安全的 map

Go 1.9+ 提供了 `sync.Map`，适合读多写少且 key 稳定的场景：

```go
var m sync.Map
m.Store("key", "value")
val, ok := m.Load("key")
m.Delete("key")
```

**注意**：`sync.Map` 没有 `Range` 方法的原子快照，遍历时数据可能变化。

## 工业界应用

### 场景 1：缓存系统

```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]*Entry
}

func (c *Cache) Get(key string) *Entry {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key string, entry *Entry) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = entry
}
```

**为什么用 RWMutex**：缓存通常是"读远多于写"，RWMutex 允许多个读请求并发，显著提升吞吐量。

### 场景 2：请求计数器

```go
type Metrics struct {
    requests atomic.Int64 // Go 1.19+
    errors   atomic.Int64
}

func (m *Metrics) IncRequests() {
    m.requests.Add(1)
}

func (m *Metrics) IncErrors() {
    m.errors.Add(1)
}

func (m *Metrics) Report() {
    fmt.Printf("requests=%d errors=%d\n", 
        m.requests.Load(), m.errors.Load())
}
```

**为什么用 atomic**：计数器频繁更新，atomic 比 Mutex 快 10 倍，且代码更简洁。

### 场景 3：连接池

```go
type ConnectionPool struct {
    mu       sync.Mutex
    conns    []*Connection
    maxConns int
}

func (p *ConnectionPool) Acquire() *Connection {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if len(p.conns) == 0 {
        return newConnection()
    }
    
    conn := p.conns[len(p.conns)-1]
    p.conns = p.conns[:len(p.conns)-1]
    return conn
}

func (p *ConnectionPool) Release(conn *Connection) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if len(p.conns) < p.maxConns {
        p.conns = append(p.conns, conn)
    } else {
        conn.Close()
    }
}
```

**为什么用 Mutex**：连接池涉及复杂逻辑（判断、切片操作），atomic 无法处理。

## 小结

### 核心要点

1. **Mutex 保护临界区**：`Lock()` 和 `Unlock()` 必须配对，推荐用 `defer`
2. **RWMutex 优化读多写少**：读锁可并发，写锁独占
3. **atomic 用于简单计数**：CPU 指令级别，性能最优
4. **用 -race 检测竞态**：CI 流程中集成竞态检测
5. **选择工具看场景**：共享状态→Mutex，消息传递→channel，计数→atomic

### 关键术语

| 英文 | 中文 | 说明 |
|------|------|------|
| race condition | 竞态条件 | 并发访问共享变量导致的不确定性 |
| critical section | 临界区 | 需要互斥访问的代码段 |
| deadlock | 死锁 | 两个 goroutine 互相等待对方释放锁 |
| atomic operation | 原子操作 | 不可分割的操作，要么全做要么全不做 |
| mutex | 互斥锁 | 同一时刻只允许一个 goroutine 持有的锁 |

### 下一步建议

1. 用 `go test -race` 扫描你的项目，修复所有竞态告警
2. 阅读 `sync` 包源码，理解 Mutex 的状态机设计
3. 学习 `golang.org/x/sync/singleflight`，解决"缓存击穿"问题

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 互斥锁 | Mutex | sync.Mutex 提供的排他锁，用于保护临界区 |
| 读写锁 | RWMutex | sync.RWMutex，允许多个读者或一个写者 |
| 原子操作 | Atomic Operation | sync/atomic 提供的 CPU 级别原子指令 |
| 竞态检测器 | Race Detector | go test -race 用于发现数据竞争的工具 |
| 临界区 | Critical Section | 同一时刻只能被一个 goroutine 执行的代码段 |
| 死锁 | Deadlock | 多个 goroutine 循环等待导致的程序卡死 |
| 饿死 | Starvation | goroutine 长期无法获得锁的情况 |
| 读多写少 | Read-Heavy | 适合使用 RWMutex 的场景 |
| 自旋 | Spinning | Mutex 在阻塞前短暂循环等待的优化策略 |
| 条件变量 | Condition Variable | sync.Cond，用于 goroutine 之间通知的机制 |

## 源码

完整示例代码位于：[internal/advance/concurrency_advanced/concurrency_advanced.go](https://github.com/savechina/hello-go/blob/main/internal/advance/concurrency_advanced/concurrency_advanced.go)
