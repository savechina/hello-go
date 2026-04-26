# 错误处理 (Error Handling)

## 开篇故事

想象你在医院看病。护士问你："哪里不舒服？"

**糟糕的回答**："不舒服。"（这相当于 `errors.New("error")`）

**有帮助的回答**："肚子痛，在右下腹，持续 2 小时，疼痛等级 7/10。"（这相当于带字段的自定义错误）

Go 的错误处理也是这样：基础阶段的 `if err != nil { return err }` 就像说"出错了"——没错，但信息太少。生产环境需要**结构化错误**：哪里出的错（operation）、哪个字段有问题（field）、输入值是什么（value）、根本原因是什么（underlying error）。

这一章教你如何设计可诊断、可分类、可追踪的错误系统，让错误从"麻烦"变成"诊断工具"。

## 本章适合谁

- ✅ 写过 `return errors.New("something went wrong")`，现在想让错误更有信息量
- ✅ 用过 `fmt.Errorf` 但不清楚 `%w` 和 `%v` 的区别
- ✅ 需要向 API 调用方返回结构化错误信息（如"字段 X 无效"）
- ✅ 想理解 `errors.Is` 和 `errors.As` 的实际应用场景

如果你曾经在日志里看到"error: failed to process request"却不知道从何查起，本章必读。

## 你会学到什么

完成本章后，你将能够：

1. **定义自定义错误类型**：实现 `Error()` 方法，携带业务相关字段
2. **使用 errors.Is 判断错误类型**：识别哨兵错误（sentinel error），判断是否权限不足、配置缺失等
3. **使用 errors.As 提取错误信息**：从错误链中提取结构化信息（字段名、输入值）
4. **用 %w 包装错误**：构建错误上下文链，便于追踪问题根源
5. **设计错误处理策略**：区分可恢复错误和致命错误，实现优雅降级

## 前置要求

在开始之前，请确保你已掌握：

- Go 接口（interface）和实现（`Error() string` 方法）
- 结构体（struct）定义和字段访问
- `fmt.Errorf` 的基本用法（`%v`、`%s` 格式化）
- 错误返回模式（`func() error`）

了解 Go 1.13+ 的错误处理特性有帮助，但本章会从基础讲起。

## 第一个例子

让我们从一个最简单的自定义错误开始：验证用户年龄。

```go
package main

import (
    "errors"
    "fmt"
)

// 自定义错误类型
type ValidationError struct {
    Field string
    Value any
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s=%v: %s", e.Field, e.Value, e.Msg)
}

// 验证函数
func validateAge(age int) error {
    if age >= 18 {
        return nil
    }
    return &ValidationError{
        Field: "age",
        Value: age,
        Msg:   "must be at least 18",
    }
}

func main() {
    err := validateAge(16)
    if err != nil {
        // 方式 1：直接打印
        fmt.Println(err) // age=16: must be at least 18
        
        // 方式 2：提取结构化信息
        var ve *ValidationError
        if errors.As(err, &ve) {
            fmt.Printf("字段=%s, 值=%v\n", ve.Field, ve.Value)
        }
    }
}
```

**运行结果**：
```
age=16: must be at least 18
字段=age, 值=16
```

**关键点**：
- 自定义错误类型实现 `Error() string` 方法，符合 Go 的 `error` 接口
- `errors.As` 可以提取具体类型，获取结构化信息
- 错误不仅是"消息"，还携带了**字段名**和**输入值**

## 原理解析

### 1. 什么是哨兵错误 (Sentinel Error)？

**定义**：哨兵错误是预先声明的、表示特定语义的错误值。

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrPermission   = errors.New("permission denied")
    ErrConfigMissing = errors.New("config missing")
)
```

**为什么需要哨兵错误**？

```go
// ❌ 错误方式：字符串比较
err := doSomething()
if err.Error() == "permission denied" { // 脆弱！错误文案可能变化
    // ...
}

// ✅ 正确方式：哨兵错误 + errors.Is
err := doSomething()
if errors.Is(err, ErrPermission) { // 稳定！比较的是变量地址
    // ...
}
```

**核心优势**：
- **稳定性**：错误文案可能变化，但变量引用不变
- **可包装**：即使被 `%w` 包装多层，`errors.Is` 仍能识别
- **可读性**：`ErrPermission` 比 `"permission denied"` 更清晰

### 2. errors.Is 的工作原理

`errors.Is` 会遍历整个错误链，逐层比较：

```go
// 错误链：fmt.Errorf("service: %w", fmt.Errorf("auth: %w", ErrPermission))

errors.Is(err, ErrPermission) // true

// 内部流程：
// 1. 比较 err == ErrPermission? → false
// 2. 调用 Unwrap() 获取下一层
// 3. 比较下一层 == ErrPermission? → false
// 4. 再调用 Unwrap() 获取再下一层
// 5. 比较再下一层 == ErrPermission? → true ✓
```

**代码实现理解**：
```go
func Is(err, target error) bool {
    for {
        if err == target {
            return true
        }
        
        // 获取下一层
        wrapper, ok := err.(interface{ Unwrap() error })
        if !ok {
            return false
        }
        err = wrapper.Unwrap()
    }
}
```

**关键点**：`%w` 会自动实现 `Unwrap()` 方法，形成错误链。

### 3. errors.As 的类型提取

`errors.As` 用于从错误链中提取具体类型：

```go
type ValidationError struct {
    Field string
    Value any
    Err   error
}

func validateUser(age int) error {
    if age < 18 {
        return &ValidationError{
            Field: "age",
            Value: age,
            Err:   errors.New("must be at least 18"),
        }
    }
    return nil
}

// 使用
err := validateUser(16)
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Printf("字段=%s, 值=%v\n", ve.Field, ve.Value)
}
```

**注意事项**：
- 第二个参数必须是**指针的指针**（`&ve`，类型是 `**ValidationError`）
- `errors.As` 也会遍历错误链，即使类型在中间某层

**为什么是指针的指针**？因为 `errors.As` 需要修改传入的变量，让它指向找到的错误类型。

### 4. %w 包装 vs %v 格式化

**%w（包装）**：
```go
err := fmt.Errorf("load config: %w", ErrConfigMissing)

// 保留原始错误，可被 errors.Is/As 识别
errors.Is(err, ErrConfigMissing) // true
```

**%v（格式化）**：
```go
err := fmt.Errorf("load config: %v", ErrConfigMissing)

// 字符串拼接，原始错误丢失
errors.Is(err, ErrConfigMissing) // false
```

**对比表**：

| 特性 | %w | %v |
|------|----|----|
| 保留原始错误 | ✅ | ❌ |
| 可被 errors.Is 识别 | ✅ | ❌ |
| 可被 errors.As 提取 | ✅ | ❌ |
| 打印时显示链 | ✅ | ✅ |
| 适用场景 | 需要保留上下文 | 仅日志打印 |

**规则**：
- 如果错误需要返回给调用方处理 → 用 `%w`
- 如果错误只用于日志打印 → 用 `%v`

### 5. 错误链的构建

通过多层包装，可以构建清晰的错误上下文链：

```go
func runJob(name string) error {
    if err := flushReport(); err != nil {
        return fmt.Errorf("run %s: %w", name, err)
    }
    return nil
}

func flushReport() error {
    if err := writeCache(); err != nil {
        return fmt.Errorf("flush report: %w", err)
    }
    return nil
}

func writeCache() error {
    if err := writeDisk(); err != nil {
        return fmt.Errorf("write cache: %w", err)
    }
    return nil
}

func writeDisk() error {
    return errors.New("disk full")
}
```

**最终错误链**：
```
run nightly-job → flush report → write cache → disk full
```

**打印结果**：
```
run nightly-job: flush report: write cache: disk full
```

**价值**：
- 知道问题发生在哪个业务流程
- 保留根本原因（disk full）
- 便于日志聚合和分析

## 常见错误

### 错误 1：用 %v 代替 %w

```go
// ❌ 错误代码
err := fmt.Errorf("load config: %v", ErrConfigMissing)

// 后果：错误链断开，无法用 errors.Is 判断
if errors.Is(err, ErrConfigMissing) { // false
    // 永远不会执行到这里
}
```

**修复**：
```go
// ✅ 修复
err := fmt.Errorf("load config: %w", ErrConfigMissing)

if errors.Is(err, ErrConfigMissing) { // true
    // 正确处理
}
```

**规则**：需要保留原始错误语义时，必须用 `%w`。

### 错误 2：errors.As 参数传错

```go
// ❌ 错误代码
var ve ValidationError // 不是指针
if errors.As(err, ve) { // 编译错误

// ❌ 错误代码（另一种）
var ve *ValidationError
if errors.As(err, ve) { // 编译错误，应该传 &ve
```

**修复**：
```go
// ✅ 正确：传指针的指针
var ve *ValidationError
if errors.As(err, &ve) {
    // ...
}
```

**记忆技巧**：`errors.As` 需要**修改变量**，所以必须传地址。

### 错误 3：自定义错误忘记实现 Unwrap()

```go
// ❌ 错误代码
type validationError struct {
    Field string
    Err   error // 底层错误
}

func (e *validationError) Error() string {
    return fmt.Sprintf("%s: %v", e.Field, e.Err)
}
// 忘记实现 Unwrap()

// 后果：errors.Is/As 无法穿透这层
err := &validationError{Err: ErrPermission}
errors.Is(err, ErrPermission) // false（错误！应该是 true）
```

**修复**：
```go
// ✅ 修复：实现 Unwrap()
func (e *validationError) Unwrap() error {
    return e.Err
}

errors.Is(err, ErrPermission) // true ✓
```

## 动手练习

### 练习 1：预测输出

阅读以下代码，预测输出（先自己想，再看答案）：

```go
var ErrNotFound = errors.New("not found")

err := fmt.Errorf("get user: %w", fmt.Errorf("query db: %w", ErrNotFound))

fmt.Println(errors.Is(err, ErrNotFound)) // ?
fmt.Println(err) // ?
```

<details>
<summary>点击查看答案</summary>

**输出**：
```
true
get user: query db: not found
```

**解析**：
- `%w` 包装后，`errors.Is` 能穿透所有层找到 `ErrNotFound`
- 打印时会显示完整错误链
</details>

### 练习 2：修复错误类型提取

以下代码为什么提取不到 `ValidationError`？如何修复？

```go
type ValidationError struct {
    Field string
    Value any
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s=%v", e.Field, e.Value)
}

func validate(age int) error {
    if age < 18 {
        return fmt.Errorf("validate: %v", &ValidationError{
            Field: "age", Value: age,
        }) // %v 断了错误链
    }
    return nil
}

// 使用
err := validate(16)
var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Field) // 不会执行
}
```

<details>
<summary>点击查看答案</summary>

**问题**：用 `%v` 包装，错误链断开。

**修复**：
```go
func validate(age int) error {
    if age < 18 {
        return fmt.Errorf("validate: %w", &ValidationError{
            Field: "age", Value: age,
        }) // 改为 %w
    }
    return nil
}
```

**或者**（不需要包装时）：
```go
func validate(age int) error {
    if age < 18 {
        return &ValidationError{
            Field: "age", Value: age,
        } // 直接返回
    }
    return nil
}
```
</details>

### 练习 3：实现 HTTP 错误响应

设计一个函数，将错误转换为 HTTP 响应（4xx/5xx）。

<details>
<summary>点击查看答案</summary>

```go
type HTTPError struct {
    Code   int
    Public string
    Err    error
}

func (e *HTTPError) Error() string {
    return e.Public
}

func (e *HTTPError) Unwrap() error {
    return e.Err
}

// 错误处理中间件
func handleError(w http.ResponseWriter, err error) {
    var httpErr *HTTPError
    
    if errors.As(err, &httpErr) {
        w.WriteHeader(httpErr.Code)
        json.NewEncoder(w).Encode(map[string]string{
            "error": httpErr.Public,
        })
        return
    }
    
    // 默认 500
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "internal server error",
    })
}

// 使用
func handler(w http.ResponseWriter, r *http.Request) {
    err := doSomething()
    if err != nil {
        handleError(w, &HTTPError{
            Code:   http.StatusBadRequest,
            Public: "invalid input",
            Err:    err,
        })
        return
    }
}
```
</details>

## 故障排查 (FAQ)

### Q1: 如何判断应该用 errors.Is 还是 errors.As？

**判断流程**：

```
需要判断错误是否属于某个类别？
  ↓
  是 → 用 errors.Is (配合哨兵错误)
  ↓
  否
  ↓
需要提取错误中的结构化信息？
  ↓
  是 → 用 errors.As (配合自定义类型)
```

**例子**：
```go
// errors.Is: 判断是否是权限错误
if errors.Is(err, ErrPermission) {
    return http.StatusForbidden
}

// errors.As: 提取验证失败详情
var ve *ValidationError
if errors.As(err, &ve) {
    log.Printf("field=%s value=%v", ve.Field, ve.Value)
}
```

### Q2: 什么时候用哨兵错误，什么时候用自定义类型？

**哨兵错误适用场景**：
- 错误语义简单（不存在、权限不足、超时）
- 只需要判断"是不是这个错误"
- 不需要携带额外信息

**自定义类型适用场景**：
- 错误需要携带字段（验证失败的字段名、输入值）
- 错误需要携带操作名（哪个业务操作失败）
- 需要区分同一类错误的不同实例

**对比**：
```go
// 哨兵错误
var ErrNotFound = errors.New("not found")

// 自定义类型
type NotFoundError struct {
    ResourceType string
    ResourceID   string
}
```

### Q3: 如何用 errors.Join 合并多个错误？

Go 1.20+ 支持 `errors.Join` 合并多个错误：

```go
func cleanup() error {
    var errs []error
    
    if err := closeFile(); err != nil {
        errs = append(errs, err)
    }
    if err := closeDB(); err != nil {
        errs = append(errs, err)
    }
    if err := closeCache(); err != nil {
        errs = append(errs, err)
    }
    
    return errors.Join(errs...) // Go 1.20+
}

// 判断
if err := cleanup(); err != nil {
    if errors.Is(err, ErrDBClosed) {
        // 可以判断是否包含某个具体错误
    }
}
```

## 知识扩展 (选学)

### 错误包装的最佳实践

**1. 添加上下文，不要重复**：
```go
// ❌ 错误：重复信息
err := fmt.Errorf("user not found: %w", ErrNotFound)

// ✅ 正确：添加上下文
err := fmt.Errorf("get user %d: %w", userID, ErrNotFound)
```

**2. 包装层数不宜过多**：
```go
// ❌ 过度包装
return fmt.Errorf("service: %w", 
    fmt.Errorf("handler: %w", 
        fmt.Errorf("db: %w", 
            fmt.Errorf("query: %w", err))))

// ✅ 适度：在边界处包装
// db 包内部：不包装
// service 层：fmt.Errorf("get user: %w", dbErr)
```

**3. 最底层错误要有意义**：
```go
// ❌ 底层是通用错误
return errors.New("error occurred")

// ✅ 底层是具体错误
return errors.New("disk full")
```

### 第三方错误处理库

**pkg/errors**（已废弃，但仍在广泛使用）：
```go
import "github.com/pkg/errors"

// 自动记录调用栈
err := errors.Wrap(doSomething(), "context")
errors.Cause(err) // 获取原始错误
```

**go.uber.org/multierr**：
```go
import "go.uber.org/multierr"

// 在函数返回前追加错误
multierr.Append(&err, closeFile())
multierr.Append(&err, closeDB())
```

**何时使用**：
- 需要调用栈 → `pkg/errors`（或 Go 1.20+ 的 `runtime.Callers`）
- 需要合并多个错误 → `errors.Join`（Go 1.20+）或 `multierr`

### 日志和监控集成

```go
// 结构化日志
logger.Error("操作失败",
    "error", err,
    "operation", "create_user",
    "user_id", userID,
)

// Sentry 错误上报
sentry.CaptureException(err)

// Prometheus 指标
errorCounter.WithLabelValues("permission_denied").Inc()
```

**关键**：从错误中提取结构化字段，用于日志、监控、告警。

## 工业界应用

### 场景 1：API 验证错误响应

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   any    `json:"value,omitempty"`
}

type ErrorResponse struct {
    Errors []FieldError `json:"errors"`
}

func validateUser(input UserInput) error {
    var errs []error
    
    if input.Name == "" {
        errs = append(errs, &FieldError{
            Field:   "name",
            Message: "name is required",
        })
    }
    
    if input.Age < 0 || input.Age > 150 {
        errs = append(errs, &FieldError{
            Field:   "age",
            Message: "age must be between 0 and 150",
            Value:   input.Age,
        })
    }
    
    if len(errs) > 0 {
        return errors.Join(errs...)
    }
    return nil
}
```

**价值**：前端可以根据 `field` 高亮具体输入框。

### 场景 2：数据库错误分类

```go
var (
    ErrDBNotFound     = errors.New("db: not found")
    ErrDBConflict     = errors.New("db: conflict")
    ErrDBConstraint   = errors.New("db: constraint violation")
)

func handleDBError(err error) error {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return ErrDBNotFound
    }
    if isConflictError(err) {
        return ErrDBConflict
    }
    if isConstraintError(err) {
        return ErrDBConstraint
    }
    return fmt.Errorf("db: %w", err)
}

// 使用
err := handleDBError(doDBOperation())
if errors.Is(err, ErrDBNotFound) {
    return http.StatusNotFound
}
if errors.Is(err, ErrDBConflict) {
    return http.StatusConflict
}
```

### 场景 3：重试策略决策

```go
type RetryableError struct {
    Err       error
    MaxRetries int
}

func (e *RetryableError) Error() string {
    return fmt.Sprintf("retryable: %v", e.Err)
}

func (e *RetryableError) Unwrap() error {
    return e.Err
}

// 使用
func processWithRetry(f func() error) error {
    var retryErr *RetryableError
    
    err := f()
    if errors.As(err, &retryErr) {
        // 是可重试错误，执行重试
        return retry(f, retryErr.MaxRetries)
    }
    
    // 不可重试，直接返回
    return err
}
```

**适用场景**：网络超时、数据库死锁、资源暂时不可用。

## 小结

### 核心要点

1. **自定义错误类型**：实现 `Error()` 方法，携带业务相关字段
2. **哨兵错误 + errors.Is**：判断错误是否属于某个已知类别
3. **自定义类型 + errors.As**：提取错误中的结构化信息
4. **%w 包装错误**：构建错误链，保留根本原因和上下文
5. **错误分类处理**：区分可重试、可恢复、致命错误

### 关键术语

| 英文 | 中文 | 说明 |
|------|------|------|
| sentinel error | 哨兵错误 | 预先声明的错误值，表示特定语义 |
| error wrapping | 错误包装 | 用 %w 包装错误，保留原始错误 |
| error chain | 错误链 | 通过包装形成的错误层级 |
| custom error type | 自定义错误类型 | 携带额外字段的错误结构体 |
| Unwrap | 解包 | 获取包装错误的内层错误 |

### 下一步建议

1. 审查项目中的 `errors.New`，替换为有意义的哨兵错误
2. 为验证逻辑添加自定义错误类型，携带字段信息
3. 在边界处用 `%w` 包装错误，添加业务上下文
4. 用 `errors.Is/As` 替换字符串比较
5. 设计项目的错误分类体系（权限、验证、数据库、网络）

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 哨兵错误 | Sentinel Error | 预定义的错误变量，用于标识特定错误类型 |
| 错误包装 | Error Wrapping | 使用 fmt.Errorf("%w", err) 包装错误，保留原始错误信息 |
| 错误链 | Error Chain | 通过多次包装形成的错误层级结构 |
| 自定义错误类型 | Custom Error Type | 实现 Error() 方法的结构体，可携带业务字段 |
| 解包 | Unwrap | errors.Unwrap() 或自动调用，获取包装错误的内层 |
| 错误判断 | Error Inspection | 使用 errors.Is 判断错误是否属于某个类型 |
| 类型提取 | Type Assertion | 使用 errors.As 从错误链中提取具体错误类型 |
| 结构化错误 | Structured Error | 携带字段信息的错误，可用于 API 响应或日志 |
| 可重试错误 | Retryable Error | 表示操作可以安全重试的错误类型 |
| 幂等性 | Idempotency | 错误判断或处理多次调用的结果一致 |

## 源码

完整示例代码位于：[internal/advance/errorhandling/errorhandling.go](https://github.com/savechina/hello-go/blob/main/internal/advance/errorhandling/errorhandling.go)
