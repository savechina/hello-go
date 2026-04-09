# 数据处理管道（Data Pipeline）

## 开篇故事

想象你经营一家快递分拣中心。每天有成千上万的包裹涌入，你雇了三个工人负责分拣。每个工人从仓库的传送带（jobs channel）上取包裹，处理完后放到另一个传送带（results channel）上。你需要协调他们的工作，确保所有包裹都被处理，而且下班时能准时关门，不会有工人还在干活。

Go 的并发模式就像这个分拣中心。Worker Pool（工作池）是你的工人，Channel（通道）是传送带，WaitGroup（等待组）是你的点名簿。掌握这些模式，你就能写出既高效又可靠的并发程序。

本章带你通过三个实战案例，学习 Go 并发的核心模式：Worker Pool、Graceful Shutdown、Fan-out/Fan-in。

## 本章适合谁

- ✅ 已经理解 goroutine 和 channel 基础，想学习实战模式的开发者
- ✅ 需要处理大量并发任务的后台服务工程师
- ✅ 想写出可控、可关闭的并发程序的 Go 学习者
- ✅ 对 Go 并发高级模式感兴趣的技术人员

如果你还没写过 goroutine，建议先完成并发基础章节。

## 你会学到什么

完成本章后，你将能够：

1. **实现 Worker Pool**：用固定数量的 goroutine 处理动态数量的任务
2. **实现 Graceful Shutdown**：让程序优雅退出，不丢任务不卡住
3. **理解 Fan-out/Fan-in**：任务分发和结果收集的标准模式
4. **正确使用 WaitGroup**：协调多个 goroutine 的完成时机
5. **处理 Channel 关闭**：避免死锁和 panic

## 前置要求

在开始之前，请确保你已掌握：

- goroutine 的启动和基本使用
- channel 的发送、接收、关闭
- `select` 语句的基本用法
- `sync.WaitGroup` 的基本概念

## 概念说明

### Worker Pool（工作池）

Worker Pool 是一种并发模式，它启动固定数量的 goroutine（称为 worker），这些 worker 从同一个 jobs channel 读取任务，处理后写入 results channel。

**核心思想**：goroutine 数量固定，任务数量动态。这避免了"每个任务一个 goroutine"的资源浪费，也避免了无限制创建 goroutine 导致的系统崩溃。

**类比**：就像餐厅厨房雇了固定数量的厨师，订单再多也是这几个厨师处理，不会因为订单多就无限雇人。

### Graceful Shutdown（优雅关闭）

优雅关闭指程序退出前完成所有进行中的任务，而不是粗暴地立即停止。

**核心思想**：通知所有 worker 停止接收新任务，等待它们完成当前任务后再退出。

**类比**：就像餐厅下班时，经理告诉厨师"做完当前这桌菜就下班"，而不是突然关灯把客人赶走。

### Fan-out/Fan-in（扇出扇入）

Fan-out 指多个 goroutine 从同一个 channel 读取数据（分发任务），Fan-in 指多个 goroutine 的结果汇总到一个 channel（收集结果）。

**核心思想**：任务分发并行化，结果收集集中化。

**类比**：就像快递分拣，多个人从同一个传送带取包裹（Fan-out），分拣后都放到同一个出库传送带（Fan-in）。

## 代码示例

### 示例 1：Worker Pool 模式

以下代码展示了 Worker Pool 的核心实现：

```go
// 创建 jobs 和 results channel
jobs := make(chan int, 5)
results := make(chan int, 5)

// 启动 3 个 worker
var wg sync.WaitGroup
for w := 1; w <= 3; w++ {
    wg.Add(1)
    go worker(w, jobs, results, &wg)
}

// 发送 5 个任务
for j := 1; j <= 5; j++ {
    jobs <- j
}
close(jobs)  // 关闭 jobs channel，通知 worker 没有新任务了

// 等待所有 worker 完成后关闭 results channel
go func() {
    wg.Wait()
    close(results)
}()

// 收集结果
for r := range results {
    fmt.Printf("结果: %d\n", r)
}
```

**worker 函数实现**：

```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()  // 确保退出时调用 wg.Done()
    for j := range jobs {  // jobs 关闭后自动退出循环
        fmt.Printf("Worker %d 处理任务 %d\n", id, j)
        time.Sleep(10 * time.Millisecond)  // 模拟处理耗时
        results <- j * 2
    }
}
```

**关键点**：
- `jobs <-chan int` 表示只读 channel，防止 worker 意外写入
- `results chan<- int` 表示只写 channel，防止 worker 意外读取
- `for j := range jobs` 在 jobs 关闭后自动退出
- `defer wg.Done()` 确保即使 panic 也调用 Done()

### 示例 2：Graceful Shutdown 模式

以下代码展示了如何优雅等待任务完成或超时退出：

```go
// 创建完成信号 channel
done := make(chan struct{})

// 启动长时间运行的任务
go func() {
    fmt.Println("模拟长时间运行的任务...")
    time.Sleep(100 * time.Millisecond)
    fmt.Println("任务完成，发送关闭信号")
    close(done)  // 发送完成信号
}()

// 等待完成或超时
select {
case <-done:
    fmt.Println("收到完成信号")
case <-time.After(200 * time.Millisecond):
    fmt.Println("超时，强制退出")
}
```

**关键点**：
- `close(done)` 是发送完成信号的标准方式（所有等待者都会收到）
- `select` 提供两条路径：正常完成或超时强制退出
- `time.After` 创建一个定时器 channel，到期后发送当前时间

### 示例 3：Fan-out/Fan-in 模式

以下代码展示了任务分发和结果收集：

```go
// 创建输入和输出 channel
input := make(chan int, 10)
output := make(chan int, 10)

// Fan-out: 启动 3 个 goroutine 处理输入
for i := 0; i < 3; i++ {
    go func(id int) {
        for n := range input {  // 所有 goroutine 共享 input
            output <- n * n  // 计算平方并写入 output
        }
    }(i)
}

// 发送数据到 input
go func() {
    for i := 1; i <= 5; i++ {
        input <- i
    }
    close(input)  // 发送完毕后关闭 input
}()

// Fan-in: 等待所有处理者完成后关闭 output
go func() {
    // 实际项目中需要 WaitGroup 等待所有 goroutine
    time.Sleep(50 * time.Millisecond)
    close(output)
}()

// 收集结果
for r := range output {
    fmt.Printf("结果: %d\n", r)
}
```

**关键点**：
- 多个 goroutine 从同一个 channel 读取是安全的（Go 自动处理竞争）
- 所有 goroutine 共享同一个输出 channel，需要协调关闭时机
- `for r := range output` 在 output 关闭后自动退出

## 知识点总结

### 核心模式对比

| 模式 | 用途 | 关键组件 |
|------|------|----------|
| Worker Pool | 固定 worker 数量处理动态任务 | jobs channel + results channel + WaitGroup |
| Graceful Shutdown | 优雅退出不丢任务 | done channel + select + time.After |
| Fan-out/Fan-in | 并行分发 + 集中收集 | 多 goroutine 共享 input/output channel |

### Channel 方向标注

```go
jobs <-chan int     // 只读 channel（防止意外写入）
results chan<- int  // 只写 channel（防止意外读取）
```

好处：编译器会阻止错误使用，提高代码安全性。

### Channel 关闭原则

1. **只关闭发送端 channel**：接收端不要关闭
2. **关闭后可继续读取**：已发送的数据仍可读完
3. `for range` 自动退出：channel 关闭后循环结束
4. **多次关闭会 panic**：确保只关闭一次

### WaitGroup 使用规范

```go
wg.Add(1)    // 在启动 goroutine 前调用
go func() {
    defer wg.Done()  // 在 goroutine 内 defer 调用
    // 处理任务
}()
wg.Wait()    // 等待所有 goroutine 完成
```

常见错误：在 goroutine 内调用 `Add()`，导致 `Wait()` 等待时机不对。

## 练习题/思考题

### 练习 1：理解 Worker Pool 流程

问题：在示例 1 中，为什么要用 goroutine 来执行 `wg.Wait()` 并关闭 results？直接在主 goroutine 中 `wg.Wait()` 再 `close(results)` 有什么问题？

<details>
<summary>点击查看答案</summary>

**答案**：如果直接在主 goroutine 中 `wg.Wait()`，会发生死锁。因为主 goroutine 等待所有 worker 完成，但 worker 完成后需要写入 results channel，而主 goroutine 还没有开始读取 results（它在等待）。没有读取者，worker 写入会阻塞。

**正确做法**：用一个单独的 goroutine 等待并关闭 results，主 goroutine 同时可以开始读取。

```go
// ❌ 错误：死锁
wg.Wait()           // 主 goroutine 等待 worker 完成
close(results)      // 此时 worker 已经阻塞在写入 results（没人读）
for r := range results { ... }

// ✅ 正确：并发处理
go func() {
    wg.Wait()
    close(results)
}()
for r := range results { ... }  // 主 goroutine 开始读取
```
</details>

### 练习 2：分析 Graceful Shutdown

问题：在示例 2 中，如果任务执行时间是 300ms，而超时设置是 200ms，会发生什么？输出是什么？

<details>
<summary>点击查看答案</summary>

**答案**：程序会输出"超时，强制退出"。

**解析**：`select` 同时等待两个 channel。`time.After(200ms)` 在 200ms 后触发，此时 `done` 还没关闭（任务还在执行）。`select` 选择第一个就绪的分支，即超时分支。

**实际场景**：这正是 Graceful Shutdown 的意义。如果任务耗时超出预期，程序不应该无限等待，而是有超时保护强制退出。

**修改建议**：如果希望任务必须完成，超时时间应该大于任务预估时间。
</details>

### 练习 3：实现带取消的 Worker Pool

问题：如果想在 Worker Pool 中加入取消机制（收到取消信号后停止处理新任务），应该怎么修改代码？请写出关键改动。

<details>
<summary>点击查看答案</summary>

**答案**：引入 `context.Context`，让 worker 监听取消信号。

```go
// 创建可取消的 context
ctx, cancel := context.WithCancel(context.Background())

// 修改 worker 函数，加入 ctx 监听
func worker(ctx context.Context, id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
        select {
        case <-ctx.Done():  // 收到取消信号
            fmt.Printf("Worker %d 被取消\n", id)
            return
        case j, ok := <-jobs:  // 尝试取任务
            if !ok {  // jobs channel 已关闭
                return
            }
            fmt.Printf("Worker %d 处理任务 %d\n", id, j)
            results <- j * 2
        }
    }
}

// 启动 worker 时传入 ctx
go worker(ctx, w, jobs, results, &wg)

// 发送取消信号
cancel()  // 所有 worker 会收到 ctx.Done() 信号
```

**关键改动**：
1. 创建 context
2. worker 内用 `select` 监听 `ctx.Done()`
3. 外部调用 `cancel()` 发送取消信号
</details>

## 常见错误

### 错误 1：在 goroutine 内调用 wg.Add()

```go
// ❌ 错误
for w := 1; w <= 3; w++ {
    go func() {
        wg.Add(1)  // 错误时机！
        defer wg.Done()
        // ...
    }()
}
wg.Wait()  // 可能还没 Add 就 Wait 了，提前退出
```

**正确做法**：

```go
// ✅ 正确
for w := 1; w <= 3; w++ {
    wg.Add(1)  // 在启动 goroutine 前调用
    go func() {
        defer wg.Done()
        // ...
    }()
}
wg.Wait()
```

### 错误 2：关闭 channel 后继续发送

```go
// ❌ 错误
close(jobs)
jobs <- 6  // panic: send on closed channel
```

**正确做法**：确保关闭后不再发送。通常用 `defer` 或在发送完毕后立即关闭。

### 错误 3：重复关闭 channel

```go
// ❌ 错误
close(results)
close(results)  // panic: close of closed channel
```

**正确做法**：只在一个地方关闭，通常用 WaitGroup 协调。

## 工业界应用

### 场景 1：批量数据处理

某公司需要每天处理百万条订单记录。使用 Worker Pool：

```go
// 实际场景：50 个 worker 处理百万条数据
jobs := make(chan Order, 1000)  // 缓冲 1000 条
results := make(chan ProcessedOrder, 1000)

// 启动 50 个 worker
for w := 0; w < 50; w++ {
    wg.Add(1)
    go processWorker(w, jobs, results, &wg)
}

// 批量发送任务（可能来自数据库查询）
for _, order := range orders {
    jobs <- order
}
close(jobs)

// 收集处理结果写入数据库
go func() {
    wg.Wait()
    close(results)
}()

for processed := range results {
    db.Save(processed)
}
```

**价值**：
- 控制资源使用（50 个 goroutine，不是百万个）
- 任务队列缓冲（1000 条缓冲，平滑突发流量）
- 结果集中收集（方便批量写入）

### 场景 2：Web 服务请求处理

```go
// 每个请求创建一个 worker pool 处理子任务
func handleRequest(ctx context.Context, items []Item) {
    jobs := make(chan Item, len(items))
    results := make(chan Result, len(items))
    
    // 启动 worker
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go worker(ctx, jobs, results, &wg)
    }
    
    // 发送任务
    go func() {
        for _, item := range items {
            jobs <- item
        }
        close(jobs)
    }()
    
    // 收集结果
    go func() {
        wg.Wait()
        close(results)
    }()
    
    for r := range results {
        // 处理结果
    }
}
```

**特点**：context 传递，请求取消时所有 worker 停止。

### 场景 3：日志聚合系统

```go
// Fan-out/Fan-in 用于多源日志聚合
input := make(chan LogEntry, 10000)
output := make(chan AggregatedLog, 100)

// Fan-out: 多个 goroutine 处理不同类型的日志
for i := 0; i < 5; i++ {
    go logProcessor(input, output)
}

// 多个日志源写入 input
go kafkaSource(input)
go fileSource(input)
go httpSource(input)

// Fan-in: 单一输出写入存储
for agg := range output {
    storage.Write(agg)
}
```

## 小结

### 核心要点

1. **Worker Pool**：固定 worker 数量处理动态任务，避免资源浪费
2. **Graceful Shutdown**：等待完成或超时退出，不丢任务不卡住
3. **Fan-out/Fan-in**：任务并行分发，结果集中收集
4. **Channel 方向标注**：`<-chan` 只读，`chan<-` 只写，提高安全性
5. **WaitGroup 规范**：`Add` 在启动前，`Done` 用 defer，`Wait` 等待完成

### 关键术语

| 英文 | 中文 | 说明 |
|------|------|------|
| Worker Pool | 工作池 | 固定数量的 goroutine 处理任务 |
| Graceful Shutdown | 优雅关闭 | 完成进行中任务后退出 |
| Fan-out | 扇出 | 多个 goroutine 从同一 channel 读取 |
| Fan-in | 扇入 | 多个 goroutine 写入同一 channel |
| WaitGroup | 等待组 | 协调多个 goroutine 完成的同步原语 |
| Channel Direction | Channel 方向 | 只读或只写的 channel 类型标注 |

### 下一步建议

1. 阅读 `sync` 包文档，了解 `Mutex`、`Cond` 等其他同步原语
2. 学习 `context` 包，掌握更优雅的取消和超时控制
3. 在项目中实践：为批量处理任务实现 Worker Pool

## 源码

完整示例代码位于：[internal/awesome/datapipeline/datapipeline.go](../../internal/awesome/datapipeline/datapipeline.go)

运行方式：
```bash
go run main.go awesome datapipeline
```