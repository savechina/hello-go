# Context 上下文 (Context)

## 开篇故事

想象你在一家繁忙的餐厅点餐。你下了订单（启动了一个 goroutine），厨师开始准备。但突然你接到电话有急事必须离开。这时你需要告诉服务员："取消我的订单"。

在 Go 中，`context` 包就扮演这个角色。它让你的程序能够在不需要某个操作的结果时，优雅地取消它，避免资源浪费。没有 context，那些后台运行的 goroutine 就像没人取的外卖——一直占着厨房（内存和 CPU）。

## 本章适合谁

- ✅ 已经会用 goroutine 和 channel，但发现 goroutine "收不住"的开发者
- ✅ 写过 HTTP 服务器，想了解如何正确处理请求超时的工程师
- ✅ 需要控制数据库查询、RPC 调用等可能耗时操作的后台服务开发者
- ✅ 想编写健壮、可取消的并发代码的 Go 学习者

如果你还在为"goroutine 泄漏"困惑，或者你的程序偶尔"卡住不退出"，本章就是为你准备的。

## 你会学到什么

完成本章后，你将能够：

1. **区分三种 context 创建方式**：`WithCancel`、`WithTimeout`、`WithDeadline`，并说出各自适用场景
2. **正确使用 cancel 函数**：理解为什么必须调用 `cancel()`，知道何时用 `defer`
3. **实现超时控制**：为任何耗时操作添加超时保护，防止程序无限等待
4. **识别 goroutine 泄漏**：通过代码审查发现缺少 context 取消的隐患
5. **在实际项目中应用 context**：将 context 作为函数第一个参数，贯穿调用链

## 前置要求

在开始之前，请确保你已掌握：

- Go 基础语法（变量、函数、结构体）
- goroutine 的启动方式（`go func()`）
- channel 的基本使用（发送、接收、`select` 语句）
- `time` 包的常用函数（`time.After`、`time.Sleep`）

如果对这些概念不熟悉，建议先阅读《并发基础》章节。

## 第一个例子

让我们从一个最简单的例子开始。假设你有一个后台任务，但你可能随时想取消它：

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // 创建一个可取消的 context
    ctx, cancel := context.WithCancel(context.Background())
    
    // 启动一个 goroutine，它会在被取消时停止
    go func() {
        select {
        case <-time.After(100 * time.Millisecond):
            fmt.Println("任务完成")
        case <-ctx.Done():
            fmt.Println("任务被取消")
        }
    }()
    
    // 主程序决定取消任务
    cancel()
    
    // 等待一下，让 goroutine 有机会执行
    time.Sleep(50 * time.Millisecond)
}
```

**运行结果**：
```
任务被取消
```

**关键点**：
- `context.Background()` 是所有 context 的"祖先"，通常只在 `main()` 或测试中使用
- `ctx.Done()` 是一个 channel，当 context 被取消时会关闭
- `cancel()` 必须调用，否则 goroutine 会一直等待（泄漏）

## 原理解析

### 1. Context 是什么？

`context.Context` 是一个接口，定义了四个方法：

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

**通俗理解**：Context 就像一个"信号广播器"。调用 `cancel()` 时，所有监听 `ctx.Done()` 的 goroutine 都会收到通知。

### 2. 为什么需要三种创建方式？

| 函数 | 用途 | 类比 |
|------|------|------|
| `WithCancel` | 手动取消 | 手动关水龙头 |
| `WithTimeout` | 超时自动取消 | 微波炉定时 |
| `WithDeadline` | 在特定时刻取消 | 闹钟在 8:00 响 |

**代码对比**：

```go
// 手动取消：适合"用户点击取消按钮"场景
ctx, cancel := context.WithCancel(context.Background())
// ... 稍后调用 cancel()

// 超时取消：适合"最多等 5 秒"场景
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel() // 必须用 defer 确保释放

// 截止时间：适合"在下午 5 点前完成"场景
ctx, cancel := context.WithDeadline(context.Background(), 
    time.Date(2026, 4, 6, 17, 0, 0, 0, time.Local))
defer cancel()
```

### 3. cancel() 为什么必须调用？

`cancel()` 的作用是：
1. 关闭 `ctx.Done()` channel，通知所有监听者
2. 释放内部资源（如定时器）

**不调用的后果**：

```go
// ❌ 错误示例：goroutine 泄漏
func badExample() {
    ctx, _ := context.WithTimeout(context.Background(), time.Hour)
    go func() {
        <-ctx.Done() // 永远等不到，因为没人调用 cancel()
    }()
    // goroutine 会一直存在，即使函数返回
}

// ✅ 正确示例
func goodExample() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
    defer cancel() // 确保函数退出时释放
    go func() {
        <-ctx.Done()
    }()
}
```

### 4. context 的传递链

Context 的核心用法是**沿着调用链传递**：

```go
func handleRequest(ctx context.Context) {
    // 传递给数据库查询
    user, err := queryUser(ctx, "alice")
    // 传递给 HTTP 请求
    resp, err := http.GetWithContext(ctx, url)
}

func queryUser(ctx context.Context, id string) (*User, error) {
    // 如果上层取消了，这里会立即返回
    rows, err := db.QueryContext(ctx, "SELECT * FROM users WHERE id = ?", id)
    // ...
}
```

**为什么这样设计**？这样，当一个 HTTP 请求被取消（如客户端断开连接），所有下游操作（数据库查询、RPC 调用）都会自动停止。

### 5. WithValue 的使用场景

`context.WithValue()` 用于传递**请求范围的元数据**：

```go
// 在请求入口处设置
ctx := context.WithValue(context.Background(), "traceID", "abc-123")
ctx := context.WithValue(ctx, "userID", 42)

// 在深层调用中读取
traceID := ctx.Value("traceID")
```

**⚠️ 注意事项**：
- 只传递轻量级元数据（trace ID、用户 ID），**不要传递大对象**
- 不要用 context 替代函数参数，它只用于"可选的"元数据
- key 最好用自定义类型，避免命名冲突

## 常见错误

### 错误 1：忘记调用 cancel()

```go
// ❌ 错误代码
ctx, _ := context.WithCancel(context.Background())
go func() {
    <-ctx.Done() // 永远不会触发
}()

// 编译器不会报错，但 goroutine 会泄漏
```

**如何修复**：
```go
// ✅ 修复：总是调用 cancel()
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // 或者在适当时机显式调用
```

### 错误 2：在 goroutine 内部调用 cancel() 但不处理 Done

```go
// ❌ 错误代码
ctx, cancel := context.WithCancel(context.Background())
go func() {
    time.Sleep(100 * time.Millisecond)
    cancel() // 取消了，但没人听
}()
// 主 goroutine 不检查 ctx.Done()，取消没有效果

// ✅ 修复：确保有 goroutine 监听 Done
go func() {
    select {
    case <-time.After(100 * time.Millisecond):
        fmt.Println("完成")
        cancel()
    case <-ctx.Done():
        fmt.Println("取消")
    }
}()
```

### 错误 3：用错 defer cancel() 的时机

```go
// ❌ 错误：defer 过早释放
func process() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel() // ❌ 函数返回时才取消，但 goroutine 可能还在用
    
    go func() {
        time.Sleep(10 * time.Second)
        // 这里 ctx 可能已经失效了
    }()
    
    return nil
}

// ✅ 正确：goroutine 和 ctx 生命周期一致
func process() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 如果 goroutine 是函数内部使用，defer 没问题
    result, err := doWork(ctx)
    return result, err
}
```

## 动手练习

### 练习 1：预测输出

阅读以下代码，预测输出结果（先自己想，再看答案）：

```go
ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()

select {
case <-time.After(100 * time.Millisecond):
    fmt.Println("A")
case <-ctx.Done():
    fmt.Println("B")
}
fmt.Println("C")
```

<details>
<summary>点击查看答案</summary>

**输出**：
```
B
C
```

**解析**：`ctx` 在 50ms 后超时，`ctx.Done()` 关闭，所以 `select` 走到第二个分支。`time.After(100ms)` 还没触发。
</details>

### 练习 2：修复 goroutine 泄漏

以下代码有什么隐患？如何修复？

```go
func startTask() {
    ctx, _ := context.WithCancel(context.Background())
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                // 做一些事
            }
        }
    }()
    // 这里忘了什么？
}
```

<details>
<summary>点击查看答案</summary>

**问题**：没有保存 `cancel` 函数，永远无法取消这个 goroutine。

**修复**：
```go
func startTask() context.CancelFunc {
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                // 做一些事
            }
        }
    }()
    return cancel // 返回取消函数，让调用方决定何时停止
}
```
</details>

### 练习 3：实现超时函数

编写一个函数 `fetchWithTimeout(url string, timeout time.Duration)`，它在指定时间内获取 URL，超时返回错误。

<details>
<summary>点击查看答案</summary>

```go
func fetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    return io.ReadAll(resp.Body)
}
```

**关键点**：`http.NewRequestWithContext` 和 `Do` 都会遵守 context 的超时设置。
</details>

## 故障排查 (FAQ)

### Q1: 如何判断我的程序有 goroutine 泄漏？

**症状**：
- 程序运行时间越长，内存占用越高
- 程序"卡住"，不退出
- 日志显示有 goroutine 一直在运行

**排查工具**：
```bash
# 使用 pprof 查看 goroutine 状态
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

**常见原因**：
- 忘记调用 `cancel()`
- channel 阻塞（发送时没人接收）
- `select` 没有 `default` 或 `Done()` 分支

### Q2: context 应该作为函数参数的第几个位置？

**答案**：**第一个参数**。

```go
// ✅ 标准写法
func Query(ctx context.Context, id string) (*User, error)

// ❌ 不推荐
func Query(id string, ctx context.Context) (*User, error)
```

**理由**：context 不属于业务参数，它是"元参数"，放在最前面便于识别和管理。

### Q3: 可以在多个 goroutine 中同时调用 cancel() 吗？

**答案**：**可以**，`cancel()` 是幂等的。

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    // 某些条件满足时取消
    if errorHappened {
        cancel() // 安全
    }
}()

go func() {
    // 超时后也取消
    <-time.After(time.Minute)
    cancel() // 即使已经被调用过，也不会 panic
}()
```

**但注意**：多次调用没有意义，通常只需要在一个地方调用。

## 知识扩展 (选学)

### context 的内部实现

`context` 的核心是一个链表结构。每次调用 `WithCancel`、`WithTimeout` 等，都会创建一个新节点，指向父节点。

```
context.Background()
       ↓
  WithCancel (父节点是 Background)
       ↓
  WithTimeout (父节点是 WithCancel)
```

当调用 `cancel()` 时，会递归关闭所有子节点。这就是为什么"取消信号"可以沿着调用链传递。

### 自定义 context

Go 官方**不建议**自定义 context 类型，但你可以用 `context.WithValue` 传递自定义数据：

```go
// 定义 key 类型（避免冲突）
type contextKey string
const userKey contextKey = "userID"

// 设置值
ctx := context.WithValue(context.Background(), userKey, 42)

// 读取值
if userID, ok := ctx.Value(userKey).(int); ok {
    fmt.Println(userID)
}
```

### context 和 errgroup

`golang.org/x/sync/errgroup` 内部使用了 context，提供更简洁的并发控制：

```go
g, ctx := errgroup.WithContext(context.Background())

g.Go(func() error {
    return doWork1(ctx) // 如果其他 goroutine 出错，ctx 会取消
})

g.Go(func() error {
    return doWork2(ctx)
})

return g.Wait() // 等待所有完成，或第一个错误
```

## 工业界应用

### 场景 1：HTTP 服务器处理请求

```go
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // HTTP 框架自动创建，客户端断开时取消
    
    // 你的处理逻辑
    user, err := s.db.GetUser(ctx, r.URL.Query().Get("id"))
    // 如果客户端断开，GetUser 会立即返回错误
    
    json.NewEncoder(w).Encode(user)
}
```

**为什么有效**：`http.Request` 的 context 会在客户端断开或超时时自动取消，所有使用该 context 的操作都会停止。

### 场景 2：批量数据处理

```go
func (s *Service) ProcessBatch(ctx context.Context, items []Item) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err() // 优雅地提前退出
        default:
        }
        
        if err := s.processOne(ctx, item); err != nil {
            return err
        }
    }
    return nil
}
```

**价值**：调用方可以随时取消批量处理，不会浪费资源处理不需要的数据。

### 场景 3：数据库连接池管理

```go
// 设置查询超时
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

rows, err := db.QueryContext(ctx, "SELECT * FROM large_table")
// 如果 30 秒没返回，自动取消，释放数据库连接
```

## 小结

### 核心要点

1. **Context 用于取消和超时**：它是管理 goroutine 生命周期的标准方式
2. **三种创建方式**：`WithCancel`（手动）、`WithTimeout`（相对时间）、`WithDeadline`（绝对时间）
3. **必须调用 cancel()**：否则会导致 goroutine 泄漏
4. **作为第一个参数传递**：沿着调用链贯穿整个请求生命周期
5. **WithValue 只传元数据**：不要用它传递业务数据

### 关键术语

| 英文 | 中文 | 说明 |
|------|------|------|
| context | 上下文 | 传递取消信号、超时的机制 |
| cancel | 取消 | 通知 goroutine 停止的信号 |
| goroutine leak | goroutine 泄漏 | goroutine 无法退出，占用资源 |
| deadline | 截止时间 | 任务必须在此时间前完成 |
| timeout | 超时 | 任务最多运行的时长 |

### 下一步建议

1. 阅读 `golang.org/x/sync/errgroup` 文档，学习更简洁的并发模式
2. 查看 `net/http` 包源码，观察 context 在 HTTP 服务器中的实际应用
3. 在你的项目中，为所有长时间运行的操作添加 context 支持

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 上下文 | Context | Go 标准库中用于在 goroutine 之间传递取消信号的机制 |
| 取消函数 | CancelFunc | context.WithCancel 返回的函数，用于取消上下文 |
| 超时 | Timeout | 使用 WithTimeout 设置的相对时间限制 |
| 截止时间 | Deadline | 使用 WithDeadline 设置的绝对时间点 |
| 背景上下文 | Background | 所有上下文的根节点，通常只在 main 函数中使用 |
| 值传递 | Value Propagation | 使用 WithValue 在调用链中传递请求范围的元数据 |
| 幂等性 | Idempotency | cancel() 可以安全调用多次，不会引发副作用 |
| 请求范围 | Request-Scoped | 与单个请求生命周期绑定的数据或操作 |

## 源码

完整示例代码位于：[internal/advance/context/context.go](https://github.com/savechina/hello-go/blob/main/internal/advance/context/context.go)
