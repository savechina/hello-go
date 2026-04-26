# 日志记录（Logging）

## 开篇故事

想象你在玩一个复杂的解谜游戏。玩到一半，你发现自己走错了路，但记不清是在哪个岔路口做错了决定。如果有个"游戏日志"，记录你每一步的选择和结果，回溯起来就容易多了。

程序也是一样的。代码运行到半夜，突然出问题了：用户投诉下单失败、支付金额不对、库存莫名其妙变负数……没有日志，你就像在黑暗中摸路，只能猜。有了日志，你可以看到：
- 用户什么时候下的单？
- 订单金额是多少？
- 库存变更前后的值是什么？

Go 提供了两套标准日志工具：
- **`log` 包**：简单直接，输出一行文本，适合调试和小工具
- **`log/slog`**：Go 1.21 新增的结构化日志，支持字段、级别、分组，适合生产环境

更重要的是，slog 允许你写**自定义处理器（Handler）**。你可以把日志存到内存里（方便测试）、写到文件、发到远程服务器，甚至可以控制"只记录警告以上级别"。

本章从最基础的 `log` 开始，逐步过渡到 slog 的结构化日志，最后实现一个自定义 Handler，理解日志系统的完整工作原理。

## 本章适合谁

- 你一直在用 `fmt.Println` 调试，想知道更专业的做法
- 你见过 `slog.Info()`，但不知道怎么添加自定义字段
- 你想理解日志级别（Info、Warn、Error）怎么用
- 你需要在测试中捕获日志输出，验证程序行为

如果你刚学 Go 基础语法，建议先理解 [函数](./functions.md) 和 [错误处理](./error-handling.md)；如果你要搭建生产环境的日志系统，可以继续学习 [日志高级用法](../advanced/logging-advanced.md)。

## 你会学到什么

学完本章，你将能够：

1. 使用 `log` 包输出简单的文本日志
2. 使用 `slog` 输出结构化日志，添加键值对字段
3. 理解日志级别（Debug、Info、Warn、Error）的意义和用法
4. 编写自定义 Handler，控制日志输出行为
5. 在测试中使用内存 Handler 捕获并验证日志

## 前置要求

在开始之前，你需要：

- **Go 1.21+**：`log/slog` 是 Go 1.21 引入的，本章示例基于 Go 1.24
- **理解函数和接口**：自定义 Handler 需要实现接口方法
- **理解 context**：slog 的 `Handle` 方法接收 `context.Context`
- **基础 I/O 概念**：知道缓冲区（Buffer）、标准输出（stdout）是什么

如果这些概念还不熟悉，建议先阅读：[接口](./interfaces.md)、[Context](../advanced/context.md)。

## 第一个例子

让我们从最简单的 `log` 包开始。它不需要导入复杂的依赖，适合快速调试。

### 使用 `log` 包

```go
func basicLogOutput(topic string) string {
    var buffer bytes.Buffer
    logger := log.New(&buffer, "basic ", 0)
    logger.Println("studying", topic)
    return strings.TrimSpace(buffer.String())
}
```

调用它：

```go
output := basicLogOutput("log package")
fmt.Println(output)
// 输出：basic studying log package
```

**关键点**：
- `bytes.Buffer`：内存缓冲区，适合测试（生产环境通常直接写 stdout）
- `"basic "`：日志前缀，每条日志都会加上
- `0`：标志位，0 表示不添加时间戳等额外信息

### 结构化日志入门

`log` 包只能输出文本，`slog` 可以输出结构化数据：

```go
func structuredLogOutput(orderID string, amount float64) string {
    var buffer bytes.Buffer
    logger := slog.New(slog.NewTextHandler(&buffer, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    logger.Info("order created", "order_id", orderID, "amount", amount)
    return strings.TrimSpace(buffer.String())
}
```

调用：

```go
output := structuredLogOutput("A-100", 19.9)
fmt.Println(output)
// 输出：time=... level=INFO msg="order created" order_id=A-100 amount=19.9
```

**优势**：
- **字段化**：`order_id` 和 `amount` 是独立字段，可以搜索、过滤
- **级别**：可以设置只记录 Info 及以上级别
- **标准化**：所有日志格式统一，便于机器解析

## 原理解析

### 1. `log` 包的工作原理

`log` 包非常简单，核心就是一个 `Logger` 结构体：

```go
type Logger struct {
    mu  sync.Mutex    // 并发安全锁
    prefix string     // 前缀
    flags int         // 标志位（时间戳、文件名等）
    out io.Writer     // 输出目标
}
```

`log.New` 创建一个 Logger：

```go
logger := log.New(&buffer, "basic ", 0)
```

`Println` 方法把消息写出去：

```go
func (l *Logger) Println(v ...interface{}) {
    l.Output(2, fmt.Sprintln(v...))
}
```

**适用场景**：
- 快速调试
- 命令行工具
- 不需要结构化的简单服务

### 2. slog 的四个核心概念

`slog` 比 `log` 复杂，有四个关键组件：

```
Logger（日志器）
  └─> Handler（处理器）
       ├─> Level（级别控制）
       └─> Attr（属性字段）
```

**Logger**：对外接口，你调用 `logger.Info()`、`logger.Warn()`。

**Handler**：实际干活的地方，决定日志怎么写、写到哪里。

**Level**：日志级别，从低到高：
- `LevelDebug = -4`
- `LevelInfo = 0`
- `LevelWarn = 4`
- `LevelError = 8`

**Attr**：键值对字段，如 `"order_id", "A-100"`。

### 3. 日志级别控制

级别控制让你在不同环境下记录不同详细程度的日志：

```go
func customHandlerOutput(minLevel slog.Level, module string) []string {
    levelVar := new(slog.LevelVar)
    levelVar.Set(minLevel)  // 动态设置级别
    
    handler := newMemoryHandler(levelVar)
    logger := slog.New(handler).With("module", module)
    
    logger.Info("skip info")    // 如果 minLevel 是 Warn，这行会被跳过
    logger.Warn("keep warn", "attempt", 2)
    logger.Error("keep error", "attempt", 3)
    
    return handler.Records()
}
```

调用：

```go
records := customHandlerOutput(slog.LevelWarn, "study")
for _, r := range records {
    fmt.Println(r)
}
// 输出：
// level=WARN msg="keep warn" module=study attempt=2
// level=ERROR msg="keep error" module=study attempt=3
```

`Info` 级别的日志被跳过了，因为 Handler 的级别设为 `Warn`。

### 4. 自定义 Handler 的实现

自定义 Handler 需要实现 `slog.Handler` 接口：

```go
type Handler interface {
    Enabled(ctx context.Context, level Level) bool
    Handle(ctx context.Context, record Record) error
    WithAttrs(attrs []Attr) Handler
    WithGroup(name string) Handler
}
```

**内存 Handler 示例**：

```go
type handlerState struct {
    records []string
}

type memoryHandler struct {
    level slog.Leveler
    attrs []slog.Attr
    group string
    state *handlerState
}

func newMemoryHandler(level slog.Leveler) *memoryHandler {
    return &memoryHandler{level: level, state: &handlerState{}}
}

func (h *memoryHandler) Enabled(_ context.Context, level slog.Level) bool {
    if h.level == nil {
        return true
    }
    return level >= h.level.Level()
}

func (h *memoryHandler) Handle(_ context.Context, record slog.Record) error {
    parts := []string{
        "level=" + record.Level.String(),
        "msg=" + record.Message,
    }
    
    // 添加全局属性
    for _, attr := range h.attrs {
        parts = append(parts, formatAttr(h.group, attr))
    }
    
    // 添加本次日志的属性
    record.Attrs(func(attr slog.Attr) bool {
        parts = append(parts, formatAttr(h.group, attr))
        return true
    })
    
    h.state.records = append(h.state.records, strings.Join(parts, " "))
    return nil
}

func (h *memoryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    clone := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
    clone = append(clone, h.attrs...)
    clone = append(clone, attrs...)
    return &memoryHandler{level: h.level, attrs: clone, group: h.group, state: h.state}
}

func (h *memoryHandler) WithGroup(name string) slog.Handler {
    nextGroup := name
    if h.group != "" {
        nextGroup = h.group + "." + name
    }
    return &memoryHandler{level: h.level, attrs: h.attrs, group: nextGroup, state: h.state}
}

func formatAttr(group string, attr slog.Attr) string {
    key := attr.Key
    if group != "" {
        key = group + "." + key
    }
    return fmt.Sprintf("%s=%v", key, attr.Value.Any())
}
```

**每个方法的作用**：
- `Enabled`：判断某个级别是否需要记录
- `Handle`：处理一条日志记录
- `WithAttrs`：添加全局属性（如 `logger.With("module", "order")`）
- `WithGroup`：添加属性分组（如 `logger.WithGroup("user").Info("msg", "id", 1)` → `user.id=1`）

### 5. 日志的并发安全

生产环境的日志通常需要并发安全。`slog` 的 Logger 本身是并发安全的：

```go
// 多个 goroutine 可以共享同一个 logger
logger := slog.New(handler)

go logger.Info("from goroutine 1")
go logger.Info("from goroutine 2")
```

但**自定义 Handler 需要自己处理并发**。上面的 `memoryHandler` 没有加锁，不适合并发场景。生产环境应该这样：

```go
type safeHandler struct {
    mu    sync.Mutex
    records []string
}

func (h *safeHandler) Handle(_ context.Context, record slog.Record) error {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.records = append(h.records, format(record))
    return nil
}
```

## 常见错误

### 错误 1：在生产环境用 `fmt.Println`

```go
// 不推荐
func processOrder(order Order) {
    fmt.Println("processing order", order.ID)
    // ...
    fmt.Println("order processed")
}
```

**问题**：
- 无法控制级别（不能只输出 Error）
- 无法结构化（不好解析）
- 无法统一前缀（看不出是哪个模块）

**修复**：

```go
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func processOrder(order Order) {
    logger.Info("processing order", "order_id", order.ID)
    // ...
    logger.Info("order processed", "order_id", order.ID)
}
```

### 错误 2：忽略日志级别配置

```go
// 所有日志都输出，包括大量 Debug
logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
logger.Debug("debug 1")
logger.Debug("debug 2")
// 生产环境会被 Debug 日志淹没
```

**修复**：根据环境设置级别。

```go
func newLogger(env string) *slog.Logger {
    opts := &slog.HandlerOptions{
        Level: slog.LevelInfo,  // 生产环境只输出 Info+
    }
    
    if env == "development" {
        opts.Level = slog.LevelDebug  // 开发环境输出所有
    }
    
    return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
```

### 错误 3：在日志中记录敏感信息

```go
// 危险！可能泄露密码
logger.Info("user login", "username", username, "password", password)
```

**修复**：脱敏或完全不记录。

```go
logger.Info("user login", "username", username)
// 或者记录哈希值
logger.Info("user login", "username", username, "password_hash", hash(password))
```

## 动手练习

### 练习 1：基本 slog 使用

创建一个 Logger，输出 Info 和 Error 级别的日志，每条至少有一个字段：

```go
func practiceLogger() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    
    // 你的代码：输出至少两条日志
}
```

<details>
<summary>参考答案</summary>

```go
func practiceLogger() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    
    logger.Info("server started", "port", 8080, "env", "production")
    logger.Error("database connection failed", "retry_count", 3, "error", "timeout")
}
```

</details>

### 练习 2：自定义 JSON Handler

修改 `memoryHandler`，让它输出 JSON 格式而不是文本格式：

```
// 文本格式：level=INFO msg="hello"
// JSON 格式：{"level":"INFO","msg":"hello"}
```

<details>
<summary>参考答案</summary>

```go
type jsonHandler struct {
    level slog.Leveler
    state *handlerState
}

func (h *jsonHandler) Handle(_ context.Context, record slog.Record) error {
    data := map[string]interface{}{
        "level": record.Level.String(),
        "msg":   record.Message,
    }
    
    record.Attrs(func(attr slog.Attr) bool {
        data[attr.Key] = attr.Value.Any()
        return true
    })
    
    jsonBytes, _ := json.Marshal(data)
    h.state.records = append(h.state.records, string(jsonBytes))
    return nil
}
```

</details>

### 练习 3：日志级别过滤实验

创建一个级别为 `LevelWarn` 的 Logger，分别调用 `Debug`、`Info`、`Warn`、`Error`，观察哪些被输出：

```go
func testLevelFilter() {
    var buffer bytes.Buffer
    handler := slog.NewTextHandler(&buffer, &slog.HandlerOptions{
        Level: slog.LevelWarn,
    })
    logger := slog.New(handler)
    
    logger.Debug("debug message")
    logger.Info("info message")
    logger.Warn("warn message")
    logger.Error("error message")
    
    fmt.Println("Output:", buffer.String())
}
```

<details>
<summary>参考答案</summary>

```go
func testLevelFilter() {
    var buffer bytes.Buffer
    handler := slog.NewTextHandler(&buffer, &slog.HandlerOptions{
        Level: slog.LevelWarn,
    })
    logger := slog.New(handler)
    
    logger.Debug("debug message")  // 不会输出
    logger.Info("info message")    // 不会输出
    logger.Warn("warn message")    // 会输出
    logger.Error("error message")  // 会输出
    
    fmt.Println("Output:", buffer.String())
    // 只有 warn 和 error 两条
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么我的日志不输出？

**A**: 最常见的原因是**级别设置不对**。检查：

```go
// 如果设置成 LevelWarn，Info 和 Debug 不会输出
opts := &slog.HandlerOptions{
    Level: slog.LevelWarn,
}
```

**解决方案**：
- 开发环境调低级别：`Level: slog.LevelDebug`
- 用 `LevelVar` 动态调整级别

```go
levelVar := new(slog.LevelVar)
levelVar.Set(slog.LevelDebug)  // 可以随时改

opts := &slog.HandlerOptions{
    Level: levelVar,
}
```

### Q2: 如何让日志带时间戳？

**A**: `slog.HandlerOptions` 有 `AddSource` 和自定义格式化：

```go
opts := &slog.HandlerOptions{
    AddSource: true,  // 添加文件名和行号
}
handler := slog.NewJSONHandler(os.Stdout, opts)
logger := slog.New(handler)
```

输出会包含：
```
{"time":"2026-04-06T10:30:00Z","level":"INFO","source":{"function":"main","file":"main.go","line":10},"msg":"hello"}
```

如果需要自定义时间格式，可以用 `slog.NewTextHandler` 并添加 `time` 字段。

### Q3: 测试中如何断言日志输出？

**A**: 用内存 Handler 捕获日志，然后断言：

```go
func TestLogger(t *testing.T) {
    state := &handlerState{}
    handler := &memoryHandler{state: state}
    logger := slog.New(handler)
    
    logger.Info("test message", "key", "value")
    
    require.Len(t, state.records, 1)
    assert.Contains(t, state.records[0], "msg=test message")
    assert.Contains(t, state.records[0], "key=value")
}
```

关键是**把日志输出变成可断言的数据**。

## 知识扩展 (选学)

### 1. 日志采样（Sampling）

高并发场景下，每条错误都记录可能导致日志爆炸。采样可以限制日志量：

```go
opts := &slog.HandlerOptions{
    Level: slog.LevelError,
}
handler := slog.NewJSONHandler(os.Stdout, opts)

// 包装成采样 Handler
sampledHandler := &samplingHandler{
    Handler: handler,
    rate:    0.1,  // 只记录 10%
}

logger := slog.New(sampledHandler)
```

采样可以大幅减少日志量，但可能漏掉重要信息，要谨慎使用。

### 2. 日志上下文（Context）

可以把日志和 `context.Context` 结合，传递请求级别的字段：

```go
type contextKey string
const loggerKey = contextKey("logger")

func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := generateID()
        logger := baseLogger.With("request_id", requestID)
        ctx := context.WithValue(r.Context(), loggerKey, logger)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func handler(w http.ResponseWriter, r *http.Request) {
    logger := r.Context().Value(loggerKey).(*slog.Logger)
    logger.Info("handling request")
}
```

这样每个请求的日志都能带上 `request_id`，方便追踪。

### 3. 日志轮转（Log Rotation）

生产环境需要定期切割日志文件，避免单个文件过大：

```go
file := &lumberjack.Logger{
    Filename:   "/var/log/myapp.log",
    MaxSize:    100,  // MB
    MaxBackups: 3,
    MaxAge:     28,   // days
}

logger := slog.New(slog.NewJSONHandler(file, nil))
```

`lumberjack` 是第三方库，专门处理日志轮转。

### 4. 分布式追踪集成

日志可以和分布式追踪（如 OpenTelemetry）集成：

```go
import (
    "go.opentelemetry.io/otel/trace"
)

func logWithTrace(ctx context.Context) {
    span := trace.SpanFromContext(ctx)
    traceID := span.SpanContext().TraceID()
    
    logger := baseLogger.With("trace_id", traceID.String())
    logger.Info("processing request")
}
```

这样日志和追踪链路可以关联起来。

### 5. 结构化日志 vs 文本日志

**文本日志**（TextHandler）：
- 优点：人类可读性好，适合开发环境
- 缺点：机器解析麻烦

**JSON 日志**（JSONHandler）：
- 优点：易于机器解析，适合 ELK、Sentry 等工具
- 缺点：人类阅读不太方便

**推荐**：开发环境用文本，生产环境用 JSON。

## 工业界应用

### 场景：电商订单服务

一个典型的电商订单服务，日志可能这样设计：

```go
type OrderService struct {
    logger *slog.Logger
    db     *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true,
    }))
    
    return &OrderService{
        logger: logger.With("service", "order"),
        db:     db,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    log := s.logger.With(
        "user_id", req.UserID,
        "items", len(req.Items),
    )
    
    log.Info("creating order")
    
    // 业务逻辑
    order, err := s.saveOrder(ctx, req)
    if err != nil {
        log.Error("failed to save order", "error", err)
        return nil, err
    }
    
    log.Info("order created", "order_id", order.ID, "total", order.Total)
    return order, nil
}
```

**关键点**：
- 用 `With` 添加服务级字段（`service`）
- 每个请求添加请求级字段（`user_id`、`items`）
- 关键操作记录 Info，错误记录 Error 并带上异常信息

### 场景：微服务链路追踪

在微服务架构中，每个请求会经过多个服务。用 `request_id` 串联日志：

```go
// 网关服务
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    requestID := r.Header.Get("X-Request-ID")
    if requestID == "" {
        requestID = generateUUID()
    }
    
    log := g.logger.With("request_id", requestID, "path", r.URL.Path)
    log.Info("incoming request")
    
    // 向下游传递 request_id
    req := r.WithContext(context.WithValue(r.Context(), "request_id", requestID))
    downstream.ServeHTTP(w, req)
}

// 订单服务
func (s *OrderService) Handle(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value("request_id").(string)
    log := s.logger.With("request_id", requestID)
    log.Info("processing order")
}
```

这样，在日志系统中搜索 `request_id` 就能看到一个请求的完整链路。

### 真实案例：Kubernetes 组件日志

Kubernetes 的组件（如 kubelet、kube-apiserver）都用结构化日志：

```
I1206 10:30:00.123456   12345 kubelet.go:2000] "SyncLoop (ADD)" source="api" pod="default/my-app-abc123"
E1206 10:30:01.234567   12345 kubelet.go:2100] "Failed to pull image" err="image not found" image="myregistry.com/app:v1"
```

格式是：`[级别][时间][组件：行号] "消息" 字段=值`。

这种格式既适合人类阅读，也方便用正则提取字段。

### 真实案例：标准库 `slog` 包

Go 官方在 `slog` 包中提供了标准实现。看看核心 API：

```go
// 创建 Logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

// 记录日志
logger.Info("message", "key", "value")
logger.Warn("message", "key", "value")
logger.Error("message", "key", "value")

// 添加全局字段
logger.With("user_id", 123).Info("user action")
```

设计简洁，易于扩展。

## 小结

本章我们学习了：

1. **`log` 包**：简单的文本日志，适合调试
2. **`slog` 结构化日志**：键值对字段、级别控制、标准化格式
3. **日志级别**：Debug、Info、Warn、Error，用于过滤不同详细程度的日志
4. **自定义 Handler**：实现 `Enabled`、`Handle`、`WithAttrs`、`WithGroup` 方法
5. **工业最佳实践**：上下文、链路追踪、JSON 格式、敏感信息脱敏

关键术语：
- **结构化日志（Structured Logging）**：用键值对记录日志，便于机器解析
- **Handler**：slog 的核心接口，决定日志怎么写、写到哪里
- **日志级别（Log Level）**：控制日志详细程度的枚举值
- **Attr（Attribute）**：日志的键值对字段

下一步建议：
- 阅读 `slog` 包官方文档：https://pkg.go.dev/log/slog
- 学习日志收集工具（如 ELK、Loki）的使用
- 了解 OpenTelemetry 的日志规范

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 日志 | Logging | 记录程序运行时的信息和事件 |
| 结构化日志 | Structured Logging | 使用键值对格式的日志，便于机器解析 |
| 日志级别 | Log Level | 日志的严重程度分级（Debug、Info、Warn、Error） |
| Handler | Handler | slog 的处理器接口，决定日志输出行为 |
| Logger | Logger | 日志器，用户调用的主要接口 |
| Attr | Attribute | 日志的键值对字段，如 `"order_id", "A-100"` |
| 缓冲器 | Buffer | 内存中的临时存储区域，用于捕获日志输出 |
| 上下文 | Context | Go 的 context 包，用于传递请求范围的值 |

## 相关资源

- [`log/slog` 官方文档](https://pkg.go.dev/log/slog)
- [Go 官方博客：Introducing slog](https://go.dev/blog/slog)
- [结构化日志最佳实践](https://www.datadoghq.com/blog/structured-logging-best-practices/)
- [OpenTelemetry 日志规范](https://opentelemetry.io/docs/specs/otel/logs/)

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/logging/logging.go)
