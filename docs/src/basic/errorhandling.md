# 错误处理（Error Handling）

## 开篇故事

想象你在一家医院看病。挂号时护士告诉你："抱歉，张医生的号已经挂完了"。这不是世界末日，只是一个需要处理的**错误情况**。医生看完病开了药，药师发现："这种药和你正在吃的药有冲突"。这又是一个错误，但可以被妥善处理。最后你去缴费，刷卡时机器显示："余额不足"。这依然不是崩溃，只是一个需要 alternativ 方案的错误。

在编程中，错误处理（Error Handling）就是程序的"医疗系统"——它不是异常（exception）那种"手术失败立即死亡"的模式，而是**显式检查、逐步处理、优雅降级**的哲学。Go 把错误当作**普通值**来对待：函数返回错误，调用者检查错误，根据错误类型决定下一步行动。这种设计让控制流清晰可见，避免了"这里为什么会崩溃"的猜测游戏。

## 本章适合谁

- 已经会写基本 Go 程序，对 `if err != nil` 感到困惑的初学者
- 从 Java/Python 转来 Go，想理解"为什么不用异常"的开发者
- 想学会正确包装错误、传递上下文的工程师
- 想提高代码健壮性和可调试性的程序员

## 你会学到什么

完成本章后，你将能够：

1. **创建和返回错误**：使用 `errors.New` 定义哨兵错误，理解错误即值的设计哲学
2. **包装错误传递上下文**：用 `fmt.Errorf` 和 `%w` 添加业务语义，保留原始错误链
3. **判断错误类型**：用 `errors.Is` 检查哨兵错误，用 `errors.As` 提取结构化错误信息
4. **实现自定义错误类型**：通过实现 `Error() string` 创建带上下文的错误
5. **设计错误处理策略**：根据场景选择忽略、记录、包装、转换错误的正确方式

## 前置要求

- 已经掌握函数返回值的基本语法
- 理解结构体和方法的定义
- 了解接口的概念（`error` 本身就是接口）
- 知道什么是异常（exception）以及其他语言的错误处理方式

## 第一个例子

让我们从一个简单的金额验证开始：

```go
package main

import (
	"errors"
	"fmt"
)

var ErrAmountMustBePositive = errors.New("amount must be positive")

func validateAmount(amount int) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}
	return nil
}

func main() {
	err := validateAmount(-1)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
	// 输出：Error: amount must be positive
}
```

这个例子展示了 Go 错误处理的核心模式：**定义错误**、**返回错误**、**检查错误**。没有异常抛出，没有 try-catch，只有明确的返回值检查。

## 原理解析

### 1. error 接口：错误即值

Go 的 `error` 是一个内置接口：

```go
type error interface {
	Error() string
}
```

**任何实现了 `Error() string` 方法的类型都是错误**。这包括：

- **简单错误**：用 `errors.New` 创建
- **格式化错误**：用 `fmt.Errorf` 创建
- **自定义错误**：实现 `Error() string` 的结构体

```go
// 简单错误
err1 := errors.New("something went wrong")

// 格式化错误
err2 := fmt.Errorf("failed to connect to %s: %w", host, underlyingErr)

// 自定义错误
type ValidationError struct {
	Field string
	Value string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("field %q value %q is invalid", e.Field, e.Value)
}
```

**为什么这样设计？** 因为错误也是程序需要处理的"值"，和其他值一样可以传递、检查、转换。

### 2. errors.New：定义哨兵错误（Sentinel Error）

**哨兵错误**是预先定义的、有稳定语义的错误值：

```go
var (
	ErrNotFound      = errors.New("not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidInput  = errors.New("invalid input")
)

func GetUser(id string) (*User, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}
	// ...
}
```

**为什么用变量而不是每次创建？** 因为哨兵错误需要在多处**比较**：

```go
user, err := GetUser("")
if err == ErrInvalidInput {  // ✅ 可以比较
	// 处理特定错误
}
```

每次 `errors.New` 会创建新实例，无法用 `==` 比较。

### 3. fmt.Errorf：包装错误添加上下文

实际业务中，裸的错误信息不够用。我们需要添加**上下文（context）**：

```go
func lookupSetting(settings map[string]string, key string) (string, error) {
	value, ok := settings[key]
	if !ok {
		// 方式 1：普通格式化（不保留错误链）
		return "", fmt.Errorf("key %q not found", key)
	}
	return value, nil
}
```

但这样会丢失原始错误信息。Go 1.13+ 引入了 `%w` 包装器：

```go
func lookupSetting(settings map[string]string, key string) (string, error) {
	value, ok := settings[key]
	if !ok {
		// 方式 2：用 %w 包装（保留错误链）
		return "", fmt.Errorf("lookup %q: %w", key, ErrSettingNotFound)
	}
	return value, nil
}
```

**`%w` vs `%v` 的区别**：
- `%v`：只格式化字符串，不保留错误链
- `%w`：包装错误，可以用 `errors.Is/As` 检查

### 4. errors.Is：检查错误链

包装后的错误长这样：`lookup "timeout": setting not found`。如何判断它包含 `ErrSettingNotFound`？

```go
err := lookupSetting(map[string]string{}, "timeout")

// ❌ 错误：直接比较会失败
if err == ErrSettingNotFound {
	// 永远不会执行
}

// ✅ 正确：用 errors.Is
if errors.Is(err, ErrSettingNotFound) {
	fmt.Println("missing setting detected")
}
```

**`errors.Is` 会遍历整个错误链**，找到匹配的哨兵错误：

```
fmt.Errorf("outer: %w",          // 第 3 层
  fmt.Errorf("middle: %w",       // 第 2 层
    ErrSettingNotFound           // 第 1 层（原始错误）
  )
)
```

`errors.Is(err, ErrSettingNotFound)` 会返回 true。

### 5. errors.As：提取结构化错误信息

有时错误包含额外信息，需要提取出来：

```go
type FieldError struct {
	Field string
	Value string
	Err   error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s %q: %v", e.Field, e.Value, e.Err)
}

func (e *FieldError) Unwrap() error {
	return e.Err  // 支持错误链
}

func parseRetryCount(raw string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, &FieldError{
			Field: "retry count",
			Value: raw,
			Err:   err,
		}
	}
	return value, nil
}

// 提取自定义错误
_, err := parseRetryCount("abc")
var fieldErr *FieldError
if errors.As(err, &fieldErr) {
	fmt.Printf("field error on %s with value %q\n", 
		fieldErr.Field, fieldErr.Value)
}
```

**`errors.As` 的作用**：遍历错误链，找到第一个匹配目标类型的错误，并赋值给目标变量。

### 6. 错误处理策略

真实代码中的错误处理模式：

```go
func summarizeError(err error) string {
	if err == nil {
		return "no error"
	}

	// 策略 1：检查特定哨兵错误
	if errors.Is(err, ErrSettingNotFound) {
		return "missing setting detected"
	}

	// 策略 2：提取结构化错误
	var fieldErr *FieldError
	if errors.As(err, &fieldErr) {
		return fmt.Sprintf("field error on %s", fieldErr.Field)
	}

	// 策略 3：兜底返回错误信息
	return err.Error()
}
```

**何时用哪种策略？**
- 需要**分类处理**：用 `errors.Is`
- 需要**提取信息**：用 `errors.As`
- 只需**记录日志**：直接用 `err.Error()`

## 常见错误

### 错误 1：忽略错误返回值

```go
// ❌ 错误：忽略错误
file, _ := os.Open("config.json")
data, _ := io.ReadAll(file)  // 如果文件没打开成功，这里会 panic

// ✅ 正确：检查每个错误
file, err := os.Open("config.json")
if err != nil {
	return fmt.Errorf("open config: %w", err)
}
data, err := io.ReadAll(file)
if err != nil {
	return fmt.Errorf("read config: %w", err)
}
```

**原则**：永远不要裸用 `_` 忽略错误，除非你 100% 确定不会失败（如 `strings.Builder` 的 Write）。

### 错误 2：只比较错误字符串

```go
// ❌ 错误：字符串比较脆弱
if err.Error() == "not found" {
	// 重构时容易破坏
}

// ✅ 正确：用 errors.Is
if errors.Is(err, ErrNotFound) {
	// 安全、可重构
}
```

**为什么？** 字符串是实现的细节，哨兵错误是稳定的契约。

### 错误 3：忘记实现 Unwrap() 导致错误链断裂

```go
type CustomError struct {
	Message string
	Err     error
}

func (e *CustomError) Error() string {
	return e.Message
}

// ❌ 错误：没有 Unwrap()，errors.Is/As 无法穿透
// ✅ 正确：添加 Unwrap()
func (e *CustomError) Unwrap() error {
	return e.Err
}
```

**Go 1.13+ 约定**：如果错误包装了另一个错误，实现 `Unwrap() error` 方法。

## 动手练习

### 练习 1：预测输出结果

```go
var ErrDB = errors.New("database error")

func query() error {
	return fmt.Errorf("query users: %w", ErrDB)
}

func main() {
	err := query()
	
	fmt.Println("err == ErrDB:", err == ErrDB)
	fmt.Println("errors.Is:", errors.Is(err, ErrDB))
}
// 问：两行输出分别是什么？
```

<details>
<summary>点击查看答案</summary>

```
err == ErrDB: false
errors.Is: true
```

**解析**：`err` 是包装后的新错误，不能用 `==` 比较。但 `errors.Is` 会遍历错误链，找到 `ErrDB`。

</details>

### 练习 2：修复错误代码

下面的代码有 4 个问题，请修复：

```go
// 问题 1：哨兵错误定义错误
var ErrInvalidID = errors.New("invalid id")  // 每次调用都创建新实例

// 问题 2：没有添加上下文
func validateID(id string) error {
	if id == "" {
		return ErrInvalidID
	}
	return nil
}

// 问题 3：忽略错误
func processID(raw string) {
	validateID(raw)  // 返回值没检查
	// ... 继续处理
}

// 问题 4：字符串比较
func handleError(err error) {
	if err.Error() == "invalid id" {
		fmt.Println("invalid ID")
	}
}
```

<details>
<summary>点击查看答案</summary>

```go
// 修复 1：用 var 定义哨兵错误
var ErrInvalidID = errors.New("invalid id")

// 修复 2：添加上下文
func validateID(id string) error {
	if id == "" {
		return fmt.Errorf("validate id: %w", ErrInvalidID)
	}
	return nil
}

// 修复 3：检查错误
func processID(raw string) error {
	if err := validateID(raw); err != nil {
		return err  // 或记录日志
	}
	// ... 继续处理
	return nil
}

// 修复 4：用 errors.Is
func handleError(err error) {
	if errors.Is(err, ErrInvalidID) {
		fmt.Println("invalid ID")
	}
}
```

</details>

### 练习 3：实现带堆栈的自定义错误

创建一个 `StackError` 类型，记录错误发生的位置：

```go
type StackError struct {
	Message  string
	FuncName string
	Line     int
	Err      error
}

// 实现 Error() string
// 实现 Unwrap() error

func main() {
	err := &StackError{
		Message:  "connection failed",
		FuncName: "connectDB",
		Line:     42,
		Err:      os.ErrNotExist,
	}
	
	fmt.Println(err.Error())
	fmt.Println(errors.Is(err, os.ErrNotExist))  // 应该输出 true
}
```

<details>
<summary>点击查看答案</summary>

```go
type StackError struct {
	Message  string
	FuncName string
	Line     int
	Err      error
}

func (e *StackError) Error() string {
	return fmt.Sprintf("%s at %s:%d: %v", 
		e.Message, e.FuncName, e.Line, e.Err)
}

func (e *StackError) Unwrap() error {
	return e.Err
}
```

**测试**：
```
connection failed at connectDB:42: file does not exist
true
```

</details>

## 故障排查 (FAQ)

### Q1: 什么时候应该返回 error，什么时候应该 panic？

**A**: 遵循以下原则：

- **返回 error**：可预见的业务错误（输入验证、网络失败、文件不存在）
- **panic**：真正的异常（逻辑 bug、违反不变量、不可恢复错误）

```go
// ✅ 返回 error
if err := db.Query(); err != nil {
	return err
}

// ✅ panic（开发阶段错误）
if user == nil {
	panic("user should never be nil here")
}
```

**经验法则**：如果错误是**预期内**的，返回 error；如果是**程序 bug**，panic。

### Q2: 如何在库代码中导出错误？

**A**: 导出哨兵错误变量，让调用方可以用 `errors.Is` 检查：

```go
// 在包 mypkg/errors.go
var ErrNotFound = errors.New("not found")

// 在包 mypkg/repo.go
func Get(id string) (*Item, error) {
	if notFound {
		return nil, ErrNotFound
	}
}

// 调用方
item, err := mypkg.Get("123")
if errors.Is(err, mypkg.ErrNotFound) {
	// 处理 404
}
```

### Q3: 错误信息应该包含什么？

**A**: 遵循"4W 原则"：

- **What**：发生了什么错误
- **Where**：在哪个操作/函数
- **Why**：根本原因（用 `%w` 包装）
- **Which**：涉及的具体数据（ID、参数值）

```go
// ❌ 信息不足
return errors.New("failed")

// ✅ 包含完整上下文
return fmt.Errorf("create user %q: %w", name, underlyingErr)
```

## 知识扩展 (选学)

### 错误分组（Error Group）

Go 1.20+ 支持 `errors.Join` 合并多个错误：

```go
func cleanup() error {
	err1 := closeFile()
	err2 := closeDB()
	err3 := closeCache()
	
	// 合并所有错误
	return errors.Join(err1, err2, err3)
}

// 检查是否包含特定错误
if errors.Is(err, err2) {
	fmt.Println("DB close failed")
}
```

### 错误格式化动词

`fmt.Errorf` 支持多个动词：

```go
// %w：包装错误（只能有一个）
fmt.Errorf("outer: %w", inner)

// %v：普通格式化
fmt.Errorf("key %q not found: %v", key, err)

// %s：字符串
fmt.Errorf("user %s not found: %s", name, reason)
```

### 第三方错误库

标准库功能有限时，可以考虑：

- **github.com/pkg/errors**：自动记录堆栈（Go 1.13+ 部分功能已内置）
- **go.uber.org/multierr**：高效的错误合并
- **github.com/rotisserie/eris**：结构化错误和堆栈

## 工业界应用

### 场景 1：HTTP API 错误响应

```go
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound     = &APIError{Code: 404, Message: "not found"}
	ErrBadRequest   = &APIError{Code: 400, Message: "bad request"}
	ErrInternal     = &APIError{Code: 500, Message: "internal error"}
)

func handleError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		w.WriteHeader(apiErr.Code)
		json.NewEncoder(w).Encode(apiErr)
	} else {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(ErrInternal)
	}
}
```

### 场景 2：数据库事务错误处理

```go
func Transfer(from, to string, amount int) error {
	return db.WithTransaction(func(tx *sql.Tx) error {
		if err := debit(tx, from, amount); err != nil {
			return fmt.Errorf("debit account %q: %w", from, err)
		}
		
		if err := credit(tx, to, amount); err != nil {
			return fmt.Errorf("credit account %q: %w", to, err)
		}
		
		return nil
	})
}

// 调用方
err := Transfer("A", "B", 100)
if errors.Is(err, ErrInsufficientBalance) {
	// 余额不足，提示用户充值
} else if err != nil {
	// 其他错误，记录日志
	log.Printf("transfer failed: %v", err)
}
```

### 场景 3：配置验证错误聚合

```go
type ConfigError struct {
	Errors []error
}

func (e *ConfigError) Error() string {
	var sb strings.Builder
	sb.WriteString("configuration validation failed:")
	for _, err := range e.Errors {
		sb.WriteString("\n  - ")
		sb.WriteString(err.Error())
	}
	return sb.String()
}

func (e *ConfigError) Unwrap() []error {
	return e.Errors
}

func validateConfig(cfg Config) error {
	var errs []error
	
	if cfg.Host == "" {
		errs = append(errs, errors.New("host is required"))
	}
	if cfg.Port <= 0 {
		errs = append(errs, errors.New("port must be positive"))
	}
	
	if len(errs) > 0 {
		return &ConfigError{Errors: errs}
	}
	return nil
}
```

## 小结

**核心要点**：
- 错误是值（error is a value），显式返回和检查
- 用 `errors.New` 定义哨兵错误，用 `fmt.Errorf` 包装上下文
- `%w` 保留错误链，`errors.Is/As` 用于检查
- 自定义错误类型实现 `Error() string` 和 `Unwrap() error`
- 永远不要忽略错误返回值

**关键术语**：
- Sentinel Error：哨兵错误，预定义的稳定错误值
- Error Wrapping：错误包装，用 `%w` 添加上下文
- Error Chain：错误链，包装错误的层次结构
- Type Assertion：类型断言，从接口提取具体类型
- Stack Trace：堆栈跟踪，记录错误发生位置

**下一步**：
- 学习 defer 和 panic/recover 机制
- 实践在项目中统一定义错误类型
- 阅读标准库 `errors`、`fmt` 包的错误处理源码

## 术语表

| 英文 | 中文 | 说明 |
|------|------|------|
| Error Handling | 错误处理 | 检查、包装、传播错误的机制 |
| Sentinel Error | 哨兵错误 | 预定义的、可比较的错误值 |
| Error Wrapping | 错误包装 | 用 %w 包装错误添加上下文 |
| Error Chain | 错误链 | 通过包装形成的错误层次结构 |
| Type Assertion | 类型断言 | 从接口提取具体类型 |
| Panic | 恐慌 | Go 的异常机制，不可恢复错误 |
| Recover | 恢复 | 从 panic 中恢复执行 |
| Stack Trace | 堆栈跟踪 | 记录函数调用链 |
| Context | 上下文 | 错误发生的环境信息 |
| Unwrap | 展开 | 获取包装的底层错误 |

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/errorhandling/errorhandling.go)
