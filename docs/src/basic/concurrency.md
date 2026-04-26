# 并发（Concurrency）

## 开篇故事

想象你在一家餐厅工作。如果只有一个厨师（单线程），他必须按顺序做菜：先切菜、再炒菜、最后摆盘。如果客人点了 10 道菜，他得一道一道做，客人会等很久。

现在你有 5 个厨师（goroutine），他们同时工作，效率大幅提升。但他们需要协调——如果两个厨师同时用同一把刀（共享内存），就会出问题。Go 的解决方案是：给每个厨师配一把刀，通过"递纸条"（channel）来沟通，而不是抢同一把刀。

这就是 Go 的并发哲学：**通过通信来共享数据，而不是通过共享数据来通信**。

---

## 本章适合谁

如果你想理解 Go 的 goroutine、channel 和并发模式，本章适合你。你需要有基本的函数和变量知识，不需要任何并发经验。

---

## 你会学到什么

完成本章后，你可以：

1. 启动 goroutine 执行异步任务
2. 使用 channel 在 goroutine 之间安全传递数据
3. 使用 `select` 处理多个 channel 或超时场景
4. 使用 `sync.WaitGroup` 等待多个 goroutine 完成
5. 识别并避免常见的并发错误（goroutine 泄漏、死锁）

---

## 前置要求

- 理解函数定义和调用
- 理解变量和类型
- 不需要任何并发经验

---

## 第一个例子

让我们从最简单的并发开始——启动一个 goroutine 并通过 channel 接收结果：

```go
ch := make(chan int)

go func() {
    ch <- 42  // 发送数据到 channel
}()

value := <-ch  // 从 channel 接收数据
fmt.Println(value)  // 输出：42
```

**关键概念**：

- `go` 关键字 - 启动一个 goroutine
- `make(chan T)` - 创建一个通道
- `<-` - 发送和接收操作符

---

## 原理解析

### 1. goroutine：轻量级执行单元

goroutine 是 Go 的并发基石。它比传统线程轻得多：

| 特征         | 线程（Thread）  | goroutine        |
| ------------ | --------------- | ---------------- |
| 初始栈大小   | 1-2 MB          | 2 KB             |
| 创建成本     | 高（系统调用）  | 低（用户态）     |
| 切换成本     | 高              | 低               |
| 单进程可创建 | 几千个          | 数十万甚至百万个 |

**为什么 goroutine 这么轻？**

- 栈是动态增长的（从 2KB 开始，按需扩展）
- 调度在用户态完成（不需要操作系统介入）
- Go 运行时（runtime）自动管理所有 goroutine

**类比**：
> 线程像重型卡车——启动慢、耗油多，但能拉很多东西。goroutine 像自行车——轻便灵活，随时可以出发。

### 2. channel：安全的通信通道

channel 是 goroutine 之间传递数据的管道。它保证：**同一时刻只有一个 goroutine 能读写**。

```go
ch := make(chan int)       // 无缓冲 channel
ch := make(chan int, 10)   // 缓冲 channel（容量 10）
```

**无缓冲 vs 缓冲**：

| 特征         | 无缓冲 channel        | 缓冲 channel          |
| ------------ | --------------------- | --------------------- |
| 发送方行为   | 阻塞直到有人接收      | 有空间时立即返回      |
| 接收方行为   | 阻塞直到有人发送      | 有数据时立即返回      |
| 适用场景     | 需要同步的场景        | 解耦发送和接收节奏    |
| 死锁风险     | 低（强制同步）        | 高（可能忘记接收）    |

**Go 的并发格言**：
> "不要通过共享内存来通信，要通过通信来共享内存。"
> — Rob Pike

### 3. select：多路复用

`select` 让你同时监听多个 channel，选择第一个准备好的：

```go
select {
case msg := <-ch1:
    fmt.Println("收到 ch1:", msg)
case msg := <-ch2:
    fmt.Println("收到 ch2:", msg)
case <-time.After(time.Second):
    fmt.Println("超时")
}
```

**类比**：
> 就像你同时等 3 个快递电话——哪个先响，你就接哪个。

### 4. sync.WaitGroup：等待多个任务完成

当你启动了多个 goroutine，需要等它们全部完成再继续：

```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Printf("Worker %d done\n", id)
    }(i)
}
wg.Wait()  // 阻塞直到所有 Done() 被调用
```

**关键规则**：
- `Add(n)` 必须在 goroutine 启动前调用
- `Done()` 必须在 goroutine 结束时调用（用 `defer` 最安全）
- `Wait()` 阻塞直到计数器归零

---

## 常见错误

### 错误 1: goroutine 泄漏（忘记接收）

```go
ch := make(chan int)
go func() {
    ch <- 42  // 发送数据
}()
// 忘记接收！goroutine 永远阻塞
```

**症状**：
- 程序卡住，不输出任何内容
- `go run` 最终报 `fatal error: all goroutines are asleep - deadlock!`

**修复方法**：

确保有人接收：
```go
ch := make(chan int)
go func() {
    ch <- 42
}()
value := <-ch  // ✅ 接收数据
fmt.Println(value)
```

---

### 错误 2: 向已关闭的 channel 发送

```go
ch := make(chan int)
close(ch)
ch <- 42  // ❌ panic!
```

**编译器不会报错，但运行时会 panic**：
```
panic: send on closed channel
```

**修复方法**：

只在确定 channel 未关闭时发送，或者用 `recover` 捕获 panic：
```go
ch := make(chan int)
go func() {
    ch <- 42  // ✅ 在关闭前发送
}()
value := <-ch
close(ch)  // 发送方负责关闭
```

**规则**：
> 只有发送方应该关闭 channel，接收方不应该关闭。

---

### 错误 3: WaitGroup 的 Add 和 Done 不匹配

```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    go func() {
        wg.Add(1)  // ❌ 错误！Add 在 goroutine 内部
        defer wg.Done()
    }()
}
wg.Wait()  // 可能先于 Add 执行
```

**症状**：
- `Wait()` 立即返回，goroutine 还没完成
- 或者 `panic: sync: negative WaitGroup counter`

**修复方法**：

`Add` 必须在启动 goroutine 之前调用：
```go
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)  // ✅ 在外部
    go func() {
        defer wg.Done()
        // do work
    }()
}
wg.Wait()
```

---

## 动手练习

### 练习 1: 预测输出

不运行代码，预测下面代码的输出：

```go
ch := make(chan string)
go func() {
    ch <- "hello"
}()
msg := <-ch
fmt.Println(msg)
```

<details>
<summary>点击查看答案</summary>

**输出**:
```
hello
```

**解析**：
1. 创建无缓冲 channel
2. 启动 goroutine 发送 "hello"
3. 主 goroutine 接收 "hello"
4. 打印并退出

</details>

---

### 练习 2: 修复死锁

下面的代码会产生死锁，请修复：

```go
func main() {
    ch := make(chan int)
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
    }()
    // 只接收一次
    fmt.Println(<-ch)
}
```

<details>
<summary>点击查看修复方法</summary>

**问题**：goroutine 发送了 5 次，但主 goroutine 只接收了 1 次。剩下的 4 次发送会永远阻塞。

**修复**：
```go
func main() {
    ch := make(chan int)
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch)  // 发送完毕，关闭 channel
    }()
    for v := range ch {  // ✅ 循环接收直到 channel 关闭
        fmt.Println(v)
    }
}
```

</details>

---

### 练习 3: 使用 WaitGroup

改写下面的代码，使用 `sync.WaitGroup` 确保所有 goroutine 完成后再打印 "all done"：

```go
func main() {
    for i := 0; i < 3; i++ {
        go func(id int) {
            fmt.Printf("Worker %d\n", id)
        }(i)
    }
    fmt.Println("all done")  // 可能先于 goroutine 执行
}
```

<details>
<summary>点击查看参考实现</summary>

```go
func main() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            fmt.Printf("Worker %d\n", id)
        }(i)
    }
    wg.Wait()
    fmt.Println("all done")
}
```

</details>

---

## 故障排查 (FAQ)

### Q: goroutine 和线程（thread）有什么区别？

**A**: 主要区别：

| 特征       | 线程              | goroutine         |
| ---------- | ----------------- | ----------------- |
| 管理方     | 操作系统          | Go 运行时         |
| 栈大小     | 固定 1-2 MB       | 动态 2 KB 起      |
| 切换成本   | 高（内核态）      | 低（用户态）      |
| 单进程数量 | 几千个            | 数十万个          |
| 通信方式   | 共享内存 + 锁     | channel 或共享内存 |

---

### Q: 什么时候用 channel，什么时候用 Mutex？

**A**: 遵循这个原则：

- **传递数据/结果** → 用 channel
- **保护共享状态** → 用 Mutex
- **不确定时** → 先用 channel（更安全）

示例：
```go
// ✅ channel：传递结果
ch := make(chan Result)
go func() { ch <- compute() }()

// ✅ Mutex：保护计数器
var mu sync.Mutex
var count int
mu.Lock()
count++
mu.Unlock()
```

---

### Q: 如何检测数据竞争（data race）？

**A**: 使用 Go 内置的 race detector：

```bash
go run -race main.go
go test -race ./...
```

它会报告所有并发读写冲突。

---

## 知识扩展 (选学)

### 缓冲 channel 的陷阱

缓冲 channel 不是队列——它只是一个有容量的管道：

```go
ch := make(chan int, 2)
ch <- 1  // ✅ 不阻塞
ch <- 2  // ✅ 不阻塞
ch <- 3  // ❌ 阻塞，直到有人接收
```

**常见误解**：
> "缓冲 channel 可以当作队列使用"

实际上，缓冲 channel 的容量只是"能容忍多少发送者不被阻塞"，它不保证顺序或持久化。

---

### Context 与并发

在生产代码中，goroutine 应该支持取消。`context.Context` 是标准方式：

```go
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

ch := make(chan Result)
go func() {
    result := compute()
    select {
    case ch <- result:
    case <-ctx.Done():
        return  // 超时，放弃结果
    }
}()
```

---

## 工业界应用：并发爬虫

**场景**：抓取多个网页，限制并发数

```go
func fetchURLs(urls []string) []string {
    var (
        wg  sync.WaitGroup
        mu  sync.Mutex
        results []string
    )

    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            content := fetch(u)
            mu.Lock()
            results = append(results, content)
            mu.Unlock()
        }(url)
    }
    wg.Wait()
    return results
}
```

**为什么这样设计**：
- 每个 URL 独立抓取（goroutine）
- 用 Mutex 保护 results 切片（共享状态）
- WaitGroup 确保所有抓取完成

---

## 小结

**核心要点**：

1. **goroutine 是轻量级的** - 创建成本低，可以轻松启动数万个
2. **channel 是安全的** - 同一时刻只有一个 goroutine 能读写
3. **select 处理多路复用** - 选择第一个准备好的 channel
4. **WaitGroup 等待完成** - Add 在外部，Done 用 defer
5. **优先用 channel 通信** - 而不是共享可变状态

**关键术语**：

- **Goroutine**: Go 的轻量级执行单元
- **Channel**: goroutine 之间的安全通信通道
- **Select**: 多路复用，监听多个 channel
- **WaitGroup**: 等待多个 goroutine 完成
- **Data Race**: 多个 goroutine 同时读写同一变量
- **Deadlock**: 所有 goroutine 都在等待，没有进展

**下一步**：

- 继续：[接口](interfaces.md)
- 回顾：[阶段复习](review-basic.md)

---

## 术语表

| English       | 中文     |
| ------------- | -------- |
| Goroutine     | goroutine（通常不翻译） |
| Channel       | 通道     |
| Concurrency   | 并发     |
| Deadlock      | 死锁     |
| Data Race     | 数据竞争 |
| Buffer        | 缓冲     |
| Mutex         | 互斥锁   |

---

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/concurrency/concurrency.go)
