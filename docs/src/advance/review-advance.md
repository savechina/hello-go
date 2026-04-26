# 阶段复习：高级进阶（Review Advance）

## 开篇故事

想象你是一位餐厅老板，刚招了一批新厨师。每个厨师都单独培训过：小王精通切菜（错误处理），小李擅长调味（反射配置），小张会炒菜（数据库操作），小赵会摆盘（Web 响应）。单独考核时，每个人都表现优秀。

但第一天营业就出了问题：客人点菜后，小王切完菜直接放在案板上（没有传递给下一个环节），小李调好味倒在地上（没有应用到菜品上），小张把菜炒焦了还说"锅的问题"（错误没有正确传播），小赵把焦菜端给客人还说"这是特色"（没有做错误转换）。

这就是很多 Go 学习者的真实写照：单独学每个知识点都能理解，但组合起来就乱套。配置校验该在哪里做？数据库错误如何传递给 HTTP 层？结构体标签到底解决了什么问题？错误应该在何处包装？

本章就是一个小型的"餐厅实战演练"。我们会构建一个完整的服务流程：从配置启动 → 请求解析 → 数据校验 → 数据库存储 → 错误处理 → HTTP 响应。走完这个闭环，你就能理解各个知识点如何协作。

## 本章适合谁

- ✅ 已完成 Go 基础章节，学过错误处理、反射、数据库、Web 的开发者
- ✅ 感觉"知识点都会但不会组合使用"的学习者
- ✅ 准备开始写真实项目的工程师
- ✅ 想理解服务启动、请求处理、错误传播整体流程的开发者

如果你还没有学习过反射、数据库或 Web 章节，建议先完成那些章节再回来。

## 你会学到什么

学完本章后，你将能够：

1. **配置校验流程**：使用反射读取结构体标签，实现配置自动校验
2. **错误边界处理**：在数据库、HTTP、业务逻辑边界正确处理和传播错误
3. **HTTP 错误映射**：将底层错误转换为合适的 HTTP 状态码和响应体
4. **完整请求链路**：理解从配置 → 请求 → 数据库 → 响应的完整数据流
5. **工程化思维**：从"会写语法"进阶到"会设计服务边界"

## 前置要求

在开始本章之前，请确保你已经掌握：

- Go 基础语法和错误处理（error handling）
- 反射基础（reflect 包，结构体标签）
- 数据库基础（GORM 或 SQL 基本操作）
- Web 基础（net/http，Handler，请求响应）
- JSON 序列化（encoding/json）

如果对上述概念还不熟悉，建议先复习相关章节。

## 第一个例子

让我们从一个最简配置校验开始，理解反射如何服务于真实场景：

```go
package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// 定义配置校验错误类型
type validationError struct {
	Problems []string
}

func (e *validationError) Error() string {
	return "validation failed: " + strings.Join(e.Problems, "; ")
}

// 配置结构体，使用标签定义规则
type reviewConfig struct {
	ServiceName string `json:"service_name" required:"true"`
	ListenPort  int    `json:"listen_port" min:"1"`
	StorageDSN  string `json:"storage_dsn" required:"true"`
}

// 通用校验函数
func validateStruct(input any) error {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	
	// 处理指针
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return &validationError{Problems: []string{"nil input"}}
		}
		val = val.Elem()
		typ = typ.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("validateStruct expects struct, got %s", val.Kind())
	}
	
	problems := make([]string, 0)
	
	// 遍历所有字段
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldTyp := typ.Field(i)
		
		// 获取字段名（优先使用 json 标签）
		fieldName := fieldTyp.Tag.Get("json")
		if fieldName == "" {
			fieldName = strings.ToLower(fieldTyp.Name)
		}
		
		// 检查 required
		if fieldTyp.Tag.Get("required") == "true" && fieldVal.IsZero() {
			problems = append(problems, fmt.Sprintf("%s is required", fieldName))
		}
		
		// 检查 min 值
		if minStr := fieldTyp.Tag.Get("min"); minStr != "" {
			minVal, _ := strconv.Atoi(minStr)
			if fieldVal.Kind() == reflect.Int && fieldVal.Int() < int64(minVal) {
				problems = append(problems, fmt.Sprintf("%s must be >= %d", fieldName, minVal))
			}
		}
	}
	
	if len(problems) == 0 {
		return nil
	}
	
	return &validationError{Problems: problems}
}

func main() {
	// 示例：校验失败的情况
	err := validateStruct(reviewConfig{
		ServiceName: "",      // 必填但为空
		ListenPort:  0,       // 小于最小值 1
		StorageDSN:  "",      // 必填但为空
	})
	
	if err != nil {
		fmt.Printf("校验失败：%v\n", err)
		// 输出：validation failed: service_name is required; listen_port must be >= 1; storage_dsn is required
	}
}
```

这个例子展示了如何用反射实现一个最小可用的配置校验器。关键点：

- 结构体标签（struct tags）承载校验规则
- 反射遍历字段，读取标签并执行校验逻辑
- 返回聚合的校验错误，而非第一个错误就返回

## 原理解析

### 概念 1：服务启动边界（Startup Boundary）

任何服务启动时都需要经历：读取配置 → 校验配置 → 初始化组件。这个过程中的每个环节都是"边界"：

```go
// 配置校验是第一个边界
cfg := reviewConfig{...}
if err := validateStruct(cfg); err != nil {
    // 配置不合法，服务不应该启动
    return fmt.Errorf("validate review config: %w", err)
}

// 第二个边界：数据库连接
db, err := gorm.Open(sqlite.Open(cfg.StorageDSN), &gorm.Config{})
if err != nil {
    // 数据库连不上，服务无法工作
    return fmt.Errorf("open review database: %w", err)
}
```

在边界处做两件事：**校验输入合法性**和**包装错误上下文**。

### 概念 2：错误包装（Error Wrapping）

Go 1.13+ 引入了 `%w` 动词来包装错误：

```go
// 底层错误
err := db.Create(&record).Error

// 包装后，保留原始错误链
return fmt.Errorf("create review course: %w", err)
```

这样做的好处：

- 上层调用者知道错误发生在"创建课程"这个动作
- 使用 `errors.Is()` 和 `errors.As()` 可以 unwrap 出底层错误
- 日志中既有上下文又有原始信息

### 概念 3：HTTP 错误边界（HTTP Error Boundary）

HTTP handler 是服务最外层的边界，负责把内部错误转换成客户端能理解的响应：

```go
func (a *reviewApp) courseHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var input courseInput
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            // 客户端错误：JSON 格式不对
            writeJSON(w, http.StatusBadRequest, map[string]any{
                "error": "invalid json payload",
            })
            return
        }
        
        if err := a.createCourse(input); err != nil {
            // 服务端错误：数据库失败等
            code := classifyStatus(err)
            writeJSON(w, code, map[string]any{
                "error": classifyError(err),
            })
            return
        }
        
        // 成功响应
        writeJSON(w, http.StatusCreated, map[string]any{
            "status":  "created",
            "service": a.config.ServiceName,
            "title":   input.Title,
        })
    })
}
```

### 概念 4：错误分类（Error Classification）

不同错误应该返回不同 HTTP 状态码：

```go
func classifyStatus(err error) int {
    var validationErr *validationError
    // 客户端错误返回 400
    if errors.Is(err, errInvalidJSON) || errors.As(err, &validationErr) {
        return http.StatusBadRequest
    }
    // 其他错误返回 500
    return http.StatusInternalServerError
}

func classifyError(err error) string {
    var validationErr *validationError
    switch {
    case errors.Is(err, errInvalidJSON):
        return errInvalidJSON.Error()  // "invalid json payload"
    case errors.As(err, &validationErr):
        return validationErr.Error()   // "validation failed: ..."
    default:
        // 不暴露内部细节给客户端
        return "internal server error"
    }
}
```

### 概念 5：数据库工作流（Database Workflow）

完整的数据库操作包含多个步骤：

```go
// 1. 启动时迁移表结构
if err := db.AutoMigrate(&courseRecord{}); err != nil {
    return nil, fmt.Errorf("migrate review database: %w", err)
}

// 2. 业务层创建记录
func (a *reviewApp) createCourse(input courseInput) error {
    // 先校验输入
    if err := validateStruct(input); err != nil {
        return fmt.Errorf("validate course input: %w", err)
    }
    
    // 再写入数据库
    record := courseRecord{Title: input.Title, Instructor: input.Instructor}
    if err := a.db.Create(&record).Error; err != nil {
        return fmt.Errorf("create review course: %w", err)
    }
    
    return nil
}

// 3. 查询统计
func (a *reviewApp) courseCount() (int64, error) {
    var count int64
    if err := a.db.Model(&courseRecord{}).Count(&count).Error; err != nil {
        return 0, fmt.Errorf("count review courses: %w", err)
    }
    return count, nil
}
```

## 常见错误

### 错误 1：在 HTTP 层暴露底层错误详情

```go
// ❌ 错误示例
if err := db.Create(&record).Error; err != nil {
    // 把 SQLite 错误细节暴露给客户端
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
}

// ✅ 正确示例
if err := db.Create(&record).Error; err != nil {
    // 记录详细错误到日志
    log.Printf("database error: %v", err)
    // 返回通用错误给客户端
    writeJSON(w, http.StatusInternalServerError, map[string]any{
        "error": "internal server error",
    })
    return
}
```

### 错误 2：忘记在边界处包装错误

```go
// ❌ 错误示例
func createCourse(input courseInput) error {
    record := courseRecord{Title: input.Title, Instructor: input.Instructor}
    return a.db.Create(&record).Error  // 调用者不知道发生了什么
}

// ✅ 正确示例
func createCourse(input courseInput) error {
    record := courseRecord{Title: input.Title, Instructor: input.Instructor}
    if err := a.db.Create(&record).Error; err != nil {
        return fmt.Errorf("create review course: %w", err)
    }
    return nil
}
```

### 错误 3：配置校验放在错误的位置

```go
// ❌ 错误示例
func (a *reviewApp) courseHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 在请求处理时校验配置（太晚了！）
        if err := validateStruct(a.config); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        // ...
    })
}

// ✅ 正确示例
func newReviewApp(cfg reviewConfig) (*reviewApp, error) {
    // 在启动时校验配置（正确时机）
    if err := validateStruct(cfg); err != nil {
        return nil, fmt.Errorf("validate review config: %w", err)
    }
    // ...
}
```

## 动手练习

### 练习 1：添加最大长度校验

扩展 `validateStruct` 函数，支持 `maxlen` 标签来限制字符串最大长度。

**提示**：参考 `min` 标签的实现，使用 `fieldVal.Type().Kind() == reflect.String` 判断。

<details>
<summary>参考答案</summary>

```go
// 检查 maxlen
if maxLenStr := fieldTyp.Tag.Get("maxlen"); maxLenStr != "" {
    maxLen, _ := strconv.Atoi(maxLenStr)
    if fieldVal.Kind() == reflect.String && fieldVal.Len() > maxLen {
        problems = append(problems, fmt.Sprintf("%s must be <= %d characters", fieldName, maxLen))
    }
}
```

</details>

### 练习 2：实现错误分类函数

编写 `classifyStatus` 和 `classifyError` 函数，区分客户端错误和服务端错误。

**提示**：使用 `errors.Is()` 和 `errors.As()` 判断错误类型。

<details>
<summary>参考答案</summary>

```go
func classifyStatus(err error) int {
    var validationErr *validationError
    if errors.Is(err, errInvalidJSON) || errors.As(err, &validationErr) {
        return http.StatusBadRequest
    }
    return http.StatusInternalServerError
}

func classifyError(err error) string {
    var validationErr *validationError
    switch {
    case errors.Is(err, errInvalidJSON):
        return errInvalidJSON.Error()
    case errors.As(err, &validationErr):
        return validationErr.Error()
    default:
        return "internal server error"
    }
}
```

</details>

### 练习 3：实现数据库计数功能

编写 `courseCount` 方法，返回数据库中课程记录的总数。

<details>
<summary>参考答案</summary>

```go
func (a *reviewApp) courseCount() (int64, error) {
    var count int64
    if err := a.db.Model(&courseRecord{}).Count(&count).Error; err != nil {
        return 0, fmt.Errorf("count review courses: %w", err)
    }
    return count, nil
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么配置校验要在启动时做，而不是在请求处理时做？

**答**：配置是服务运行的前提条件。如果配置不合法，服务根本不应该启动。在启动时校验可以：

- 快速失败（fail-fast），避免问题服务上线
- 减少运行时开销（校验只做一次）
- 明确责任边界（配置错误 vs 请求错误）

### Q2: `%w` 包装错误和 `fmt.Sprintf` 拼接错误有什么区别？

**答**：`%w` 创建了错误链（error chain），可以用 `errors.Unwrap()` 逐层 unwrap：

```go
err := fmt.Errorf("outer: %w", innerErr)
errors.Is(err, innerErr)  // true - 可以检测到包裹的底层错误

// 而 Sprintf 只是字符串拼接
err2 := fmt.Sprintf("outer: %v", innerErr)
// 无法用 errors.Is() 检测底层错误
```

### Q3: 为什么要区分 `validationError` 和普通 error？

**答**：区分错误类型便于分类处理：

- **客户端错误**（如校验失败）：返回 400，帮助客户端修正请求
- **服务端错误**（如数据库失败）：返回 500，不暴露内部细节

通过类型断言或 `errors.As()` 可以精确分类错误。

## 知识扩展 (选学)

### 扩展 1：使用验证库

生产环境常用成熟验证库如 `go-playground/validator`：

```go
import "github.com/go-playground/validator/v10"

type Config struct {
    ServiceName string `validate:"required"`
    Port        int    `validate:"required,min=1,max=65535"`
}

validate := validator.New()
err := validate.Struct(cfg)
```

### 扩展 2：错误类型层次

构建更细粒度的错误类型层次：

```go
type BadRequestError struct{ ... }
type NotFoundError struct{ ... }
type DatabaseError struct{ ... }
```

### 扩展 3：中间件错误处理

在 HTTP 中间件中统一处理错误：

```go
func errorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 统一错误处理逻辑
    })
}
```

## 工业界应用

### 场景：微服务配置管理

某公司的微服务平台管理着上百个服务，每个服务有不同的配置项。平台需要：

1. **启动校验**：服务启动前强制校验配置合法性
2. **热更新**：配置变更时重新校验再应用
3. **错误定位**：配置错误时精确定位到字段
4. **文档生成**：从结构体标签自动生成配置文档

**实现方案**：

```go
type ServiceConfig struct {
    ServiceName string   `json:"service_name" required:"true" description:"服务名称"`
    Port        int      `json:"port" required:"true" min:"1" max:"65535" description:"监听端口"`
    DSN         string   `json:"dsn" required:"true" format:"url" description:"数据库连接串"`
    LogLevel    string   `json:"log_level" default:"info" enum:"debug,info,warn,error" description:"日志级别"`
}

// 启动时校验
func loadConfig(path string) (*ServiceConfig, error) {
    data, _ := os.ReadFile(path)
    var cfg ServiceConfig
    json.Unmarshal(data, &cfg)
    
    if err := validateStruct(cfg); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    // 应用默认值
    if cfg.LogLevel == "" {
        cfg.LogLevel = "info"
    }
    
    return &cfg, nil
}
```

这种模式被广泛应用于配置中心、API 网关、服务发现等基础设施。

## 小结

本章通过一个完整的服务示例，串联了反射、错误处理、数据库、Web 四个关键知识点。

### 核心链路

```
配置加载 → 反射校验 → 数据库初始化 → HTTP Handler → 错误分类 → HTTP 响应
```

### 关键原则

1. **边界思维**：在服务边界做校验和错误转换
2. **快速失败**：配置问题在启动时暴露
3. **错误包装**：向上传递时增加上下文
4. **错误隔离**：不向客户端暴露内部细节

### 下一步

- 学习更复杂的错误处理模式（retry、circuit breaker）
- 研究成熟框架（Gin、Echo）的错误处理机制
- 实践编写完整的微服务配置系统

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 边界 | Boundary | 系统/组件/层次之间的分界点 |
| 配置校验 | Configuration Validation | 检查配置合法性的过程 |
| 错误包装 | Error Wrapping | 用 %w 创建错误链 |
| 错误分类 | Error Classification | 根据错误类型返回不同响应 |
| 结构体标签 | Struct Tag | 附加在字段上的元数据 |
| 快速失败 | Fail Fast | 尽早暴露错误的设计原则 |
| 错误链 | Error Chain | 通过 %w 链接的多层错误 |
| HTTP 状态码 | HTTP Status Code | HTTP 响应的状态标识 |
| 反射校验 | Reflection-based Validation | 基于反射的通用校验逻辑 |
| 数据库迁移 | Database Migration | 自动创建/更新表结构 |

## 源码

完整示例代码位于：[internal/advance/review/review.go](https://github.com/savechina/hello-go/blob/main/internal/advance/review/review.go)
