# 包管理（Packages）

## 开篇故事

想象你在整理一个巨大的工具箱。一开始，所有工具都堆在一起：锤子、螺丝刀、扳手、电钻……找个东西得翻半天。后来你决定给工具分类：电动工具放一个箱子，手动工具放另一个箱子，测量工具单独放小盒子。每个箱子上贴个标签，找东西时直接去对应的箱子拿。

Go 的包（Package）就是这个思想。代码多了，不能全塞在一个文件里。包帮我们：
- **组织代码**：相关功能放一起，比如 `database` 包管数据库，`http` 包管网络
- **控制可见性**：有些东西只给自己人用（包内可见），有些可以对外公开（包外可见）
- **避免命名冲突**：两个包都可以有 `Config` 结构体，用 `db.Config` 和 `http.Config` 区分

更重要的是，Go 的包系统背后还有模块（Module）和依赖管理。`go.mod` 文件定义了你的项目从哪儿开始，导入路径怎么写。理解包，不只是知道 `package` 关键字，而是要理解整个代码组织的生态系统。

## 本章适合谁

- 你写过 `import "fmt"`，但不清楚 `import "hello/internal/xxx"` 是怎么工作的
- 你见过 `init()` 函数，但不知道它什么时候执行、有什么用
- 你搞不懂为什么有的标识符首字母大写、有的小写
- 你想创建一个可复用的 Go 模块，让别人能 `import` 你的代码

如果你刚学 Go 语法，建议先理解 [函数](./functions.md) 和 [结构体](./structs.md)；如果你要发布模块或管理外部依赖，可以继续学习 [模块管理](../advanced/go-mod.md)。

## 你会学到什么

学完本章，你将能够：

1. 理解 Go 包的组织方式和导入路径规则
2. 正确使用导出（exported）和未导出（unexported）标识符
3. 掌握 `init()` 函数的执行时机和适用场景
4. 理解 `go.mod` 如何定义模块路径并影响导入
5. 设计清晰的包结构，避免循环依赖和过度暴露

## 前置要求

在开始之前，你需要：

- **Go 基础语法**：理解 `package`、`import`、`func` 这些基本概念
- **文件目录概念**：知道 Go 源码文件放在什么位置
- **理解 `go.mod`**：至少见过这个文件，知道它是 Go 模块的配置文件
- **使用过标准库**：比如 `import "fmt"`、`import "os"`

如果这些概念还不熟悉，建议先阅读：[项目结构](../getting-started/project-structure.md)、[Go 模块入门](https://go.dev/doc/modules/create)。

## 第一个例子

让我们从一个真实的项目结构开始。假设你的模块名叫 `hello`，目录结构如下：

```
hello/
├── go.mod
└── internal/
    └── basic/
        └── packages/
            ├── packages.go
            └── demo/
                ├── visibility/
                │   └── visibility.go
                ├── beta/
                │   └── beta.go
                └── trace/
                    └── trace.go
```

### 导入本地包

在 `packages.go` 中，你可以这样导入子包：

```go
import (
    "hello/internal/basic/packages/demo/beta"
    "hello/internal/basic/packages/demo/trace"
    "hello/internal/basic/packages/demo/visibility"
)
```

这里的导入路径由两部分组成：
- **模块名**：`hello`（来自 `go.mod` 的 `module hello`）
- **相对路径**：`/internal/basic/packages/demo/beta`

完整路径就是 `hello/internal/basic/packages/demo/beta`。

### 使用包中的导出内容

导入后，就可以使用包里的导出标识符：

```go
func describeImportUsage(name string, score int) string {
    profile := visibility.NewProfile(name, score)
    return fmt.Sprintf("%s | %s", beta.Description(), profile.PublicSummary())
}
```

注意：只能访问**首字母大写**的标识符（如 `NewProfile`、`Description`）。

## 原理解析

### 1. 包的可见性规则

Go 用**首字母大小写**控制可见性，这是最简单也最重要的规则：

```go
// visibility/visibility.go
package visibility

// 首字母大写 = 导出（exported），包外可以访问
type Profile struct {
    Name  string  // 导出字段
    score int     // 未导出字段，包外不能直接访问
}

// 导出函数
func NewProfile(name string, score int) Profile {
    return Profile{Name: name, score: score}
}

// 未导出方法（小写）
func (p Profile) internalNote() string {
    return "for internal use only"
}

// 导出方法（大写）
func (p Profile) PublicSummary() string {
    return fmt.Sprintf("%s (score: %d)", p.Name, p.score)
}
```

**规则总结**：
- **包级别**：首字母大写 = 导出，小写 = 未导出
- **结构体字段**：同样适用，大写可访问，小写不可访问
- **方法**：接收者类型不影响，方法名首字母决定可见性

### 2. `go.mod` 和导入路径

`go.mod` 定义了模块的根路径，决定了导入怎么写：

```go
module hello

go 1.24
```

这意味着：
- **本地导入**：`hello/internal/xxx` 指当前项目的 `internal/xxx` 目录
- **外部导入**：`github.com/gin-gonic/gin` 指向远程仓库

如果修改模块名，所有导入路径都要跟着改：

```go
module github.com/weirenyan/hello

// 导入也要改
import "github.com/weirenyan/hello/internal/xxx"
```

### 3. `init()` 函数

`init()` 是特殊的初始化函数，不需要手动调用：

```go
// packages.go
var initOrder []string

func init() {
    chapters.Register("basic", "packages", Run)
    trace.Record("main.init")
    initOrder = trace.Events()
}

// trace/trace.go
var events []string

func init() {
    trace.Record("trace.init")
}

func Record(event string) {
    events = append(events, event)
}

func Events() []string {
    result := make([]string, len(events))
    copy(result, events)
    return result
}
```

**执行顺序**：
1. 先执行导入包的 `init()`（按导入顺序）
2. 再执行当前包的 `init()`
3. 最后执行 `main()`

在上例中，如果 `packages.go` 导入了 `trace`，那么：
- `trace.init()` 先执行 → `Record("trace.init")`
- `packages.init()` 后执行 → `Record("main.init")`
- `initOrder` 最终是 `["trace.init", "main.init"]`

### 4. `internal` 包的特殊规则

Go 有一个特殊约定：放在 `internal` 目录下的包，**只能被同一模块内的代码导入**。

```go
// 这是允许的（同一模块内）
import "hello/internal/basic/packages/demo/visibility"

// 这是禁止的（其他模块想导入）
// 模块 B 的 code.go
import "hello/internal/basic/packages/demo/visibility"  // 编译错误！
```

这个规则的作用是：
- **封装内部实现**：防止外部依赖你的内部细节
- **稳定公开 API**：只有 `internal` 之外的包才是公开 API

### 5. 避免循环依赖

Go 不允许循环依赖。如果 A 导入 B，B 就不能导入 A：

```
// 错误示例
package A
import "hello/B"  // A 导入 B

package B
import "hello/A"  // B 导入 A → 编译错误！
```

**解决方法**：
- 提取公共接口到第三个包 C，让 A 和 B 都导入 C
- 重新设计架构，消除双向依赖

## 常见错误

### 错误 1：在包外访问未导出标识符

```go
// visibility/visibility.go
package visibility

type Profile struct {
    Name  string
    score int  // 小写，未导出
}

// main.go
import "hello/internal/basic/packages/demo/visibility"

func main() {
    p := visibility.NewProfile("Alice", 90)
    fmt.Println(p.Name)   // ✓ 可以
    fmt.Println(p.score)  // ✗ 编译错误：score 未导出
}
```

**修复**：通过导出方法间接访问。

```go
// visibility/visibility.go
func (p Profile) GetScore() int {
    return p.score
}

// main.go
fmt.Println(p.GetScore())  // ✓ 通过导出方法访问
```

### 错误 2：在 `init()` 中放业务逻辑

```go
// 错误做法
func init() {
    // 不应该在这里调 API、写数据库、处理业务
    db, _ := sql.Open("mysql", "...")
    db.Exec("INSERT INTO logs ...")
}
```

**修复**：`init()` 只做初始化和注册。

```go
// 正确做法
func init() {
    chapters.Register("basic", "packages", Run)
    trace.Record("init")
}

// 业务逻辑放普通函数
func Run() {
    // 这里才是业务逻辑
}
```

### 错误 3：导入路径写错

```go
// 错误：忘记加模块名前缀
import "internal/basic/packages/demo/visibility"  // ✗ 编译错误

// 错误：路径拼写错误
import "hello/internal/basic/pakcages/demo/visibility"  // ✗ 编译错误
```

**修复**：始终用完整的模块路径。

```go
// 正确
import "hello/internal/basic/packages/demo/visibility"
```

可以用 `go list -m` 查看当前模块名。

## 动手练习

### 练习 1：创建导出规则

创建一个 `calculator` 包，包含以下要求：
- 一个 `Calculator` 结构体，有未导出的 `history []int` 字段
- 导出方法 `Add(n int)`、`Subtract(n int)`、`GetHistory()`
- 未导出方法 `record(n int)` 用于内部记录历史

<details>
<summary>参考答案</summary>

```go
// calculator/calculator.go
package calculator

type Calculator struct {
    history []int
}

func (c *Calculator) Add(n int) {
    c.record(n)
}

func (c *Calculator) Subtract(n int) {
    c.record(-n)
}

func (c *Calculator) GetHistory() []int {
    result := make([]int, len(c.history))
    copy(result, c.history)
    return result
}

func (c *Calculator) record(n int) {
    c.history = append(c.history, n)
}
```

</details>

### 练习 2：包初始化顺序实验

创建三个包 `alpha`、`beta`、`gamma`，每个包里都有 `init()` 打印信息。在 `main.go` 中导入它们，观察执行顺序。

<details>
<summary>参考答案</summary>

```go
// alpha/alpha.go
package alpha
import "fmt"
func init() {
    fmt.Println("alpha init")
}

// beta/beta.go
package beta
import "fmt"
func init() {
    fmt.Println("beta init")
}

// gamma/gamma.go
package gamma
import "fmt"
func init() {
    fmt.Println("gamma init")
}

// main.go
package main
import (
    _ "hello/alpha"
    _ "hello/beta"
    _ "hello/gamma"
)
func main() {
    fmt.Println("main")
}
```

输出：
```
alpha init
beta init
gamma init
main
```

</details>

### 练习 3：使用 `internal` 封装

创建 `internal/config` 包，提供一个 `Load()` 函数。然后在主程序中导入并使用它。尝试在另一个模块（如 `../other-module`）中导入，观察编译错误。

<details>
<summary>参考答案</summary>

```go
// internal/config/config.go
package config

import "fmt"

func Load() string {
    return "config loaded"
}

// main.go
package main
import (
    "hello/internal/config"
    "fmt"
)
func main() {
    fmt.Println(config.Load())
}

// ../other-module/main.go
package main
import "hello/internal/config"  // ✗ 编译错误：use of internal package not allowed
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么我的包导入后说 "not defined"？

**A**: 最常见的原因是**标识符未导出**。检查：
- 结构体名、函数名、变量名是否首字母大写
- 拼写是否正确（Go 区分大小写）

```go
// package foo
type config struct {}  // 小写，未导出

// main.go
import "hello/foo"
var c foo.config  // ✗ 错误：config 未导出

// 修复
type Config struct {}  // 大写，导出
```

### Q2: `init()` 会被调用多次吗？

**A**: 不会。每个包的 `init()` 在整个程序生命周期中**只执行一次**，即使它被多个地方导入。

```go
// 包 A 导入 C
// 包 B 也导入 C
// C 的 init() 只执行一次，不是两次
```

如果需要在每次调用时执行初始化逻辑，用普通函数：

```go
func Initialize() {
    // 每次调用都会执行
}
```

### Q3: 如何解决 "import cycle not allowed" 错误？

**A**: 循环依赖的解决思路：

1. **提取公共接口**：把双向依赖变成单向
   ```
   原结构：A ↔ B
   新结构:  A → C ← B（C 是接口或共享类型）
   ```

2. **依赖注入**：通过参数传递，而不是直接导入
   ```go
   // 不直接导入 B
   func Process(data Data, saver Saver) {
       saver.Save(data)  // Saver 是接口
   }
   ```

3. **重新设计架构**：有些循环依赖说明设计有问题，考虑合并包或调整职责

## 知识扩展 (选学)

### 1. 包的别名导入

当包名冲突或路径太长时，可以用别名：

```go
import (
    old "hello/v1/handler"
    new "hello/v2/handler"
    
    httputil "github.com/google/go-cmp/cmp/cmpopts"
)

old.Handle()
new.Handle()
```

### 2. 空白导入 `_`

导入包但不使用它的导出标识符，通常是为了触发 `init()`：

```go
import (
    _ "github.com/go-sql-driver/mysql"  // 注册数据库驱动
    _ "hello/internal/metrics"          // 注册指标收集器
)
```

这种方式叫"副作用导入"（side-effect import）。

### 3. 点导入 `.`（不推荐）

点导入可以直接使用包中的导出标识符，省略包名前缀：

```go
import . "fmt"

Println("hello")  // 等价于 fmt.Println("hello")
```

**为什么不推荐**：
- 不清楚标识符从哪来
- 可能和本地变量名冲突
- 降低可读性

只在测试文件中偶尔使用（如 `import . "github.com/onsi/ginkgo"`）。

### 4. 相对导入 `.`（已废弃）

Go 1.20+ 已经不再支持相对导入：

```go
import "./packages"  // ✗ 不支持
import "../common"   // ✗ 不支持
```

始终使用完整模块路径。

### 5. 包的测试文件

测试文件和被测包在同一个包，但用 `_test` 后缀：

```go
// calculator/calculator.go
package calculator

// calculator/calculator_test.go
package calculator  // 同名包，可以访问未导出内容

// calculator/calculator_external_test.go
package calculator_test  // 外部测试，只能访问导出内容
```

外部测试更接近真实使用场景，推荐优先使用。

## 工业界应用

### 场景：微服务项目的包结构

某电商平台的订单服务，包结构设计如下：

```
order-service/
├── cmd/
│   └── server/
│       └── main.go          # 程序入口
├── internal/
│   ├── handler/             # HTTP 处理器
│   │   ├── order.go
│   │   └── user.go
│   ├── service/             # 业务逻辑层
│   │   ├── order.go
│   │   └── payment.go
│   ├── repository/          # 数据访问层
│   │   ├── mysql/
│   │   └── redis/
│   └── config/              # 配置管理
├── pkg/
│   └── models/              # 公开的数据模型
│       └── order.go
└── go.mod
```

**关键点**：
- `internal/`：内部实现，外部不能依赖
- `pkg/`：公开 API，其他服务可以导入
- `cmd/`：可执行程序入口

### 场景：SDK 开发

开发一个 SDK 让别人使用时，包设计更讲究：

```go
// sdk/client.go
package sdk

type Client struct {
    apiKey string
}

func NewClient(apiKey string) *Client {
    return &Client{apiKey: apiKey}
}

func (c *Client) DoRequest(ctx context.Context, req Request) (*Response, error) {
    // 实现
}

// sdk/types.go
package sdk

type Request struct {
    Method string
    Path   string
}

type Response struct {
    StatusCode int
    Body       []byte
}
```

用户使用时：

```go
import "github.com/company/sdk"

client := sdk.NewClient("api-key")
resp, _ := client.DoRequest(ctx, sdk.Request{...})
```

### 真实案例：标准库 `net/http` 包

看看 Go 标准库的组织方式：

```go
import "net/http"

// 导出类型
type Client struct {}
type Server struct {}
type Request struct {}
type ResponseWriter interface {}

// 导出函数
func Get(url string) (*Response, error)
func ListenAndServe(addr string, handler Handler) error
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
```

清晰的导出规则让用户只需要关注公开 API，内部实现细节被完全隐藏。

## 小结

本章我们学习了：

1. **导入路径**：由模块名 + 相对路径组成
2. **可见性规则**：首字母大写导出，小写未导出
3. **`init()` 函数**：自动执行，用于初始化和注册
4. **`internal` 包**：限制外部依赖，封装内部实现
5. **避免循环依赖**：合理设计包结构，用接口解耦

关键术语：
- **导出（Exported）**：首字母大写，包外可访问
- **未导出（Unexported）**：首字母小写，包内专用
- **模块路径（Module Path）**：`go.mod` 定义的导入前缀
- **循环依赖（Import Cycle）**：A 导入 B、B 导入 A，Go 不允许

下一步建议：
- 阅读 Go 官方博客 "Organizing Go Modules"
- 学习标准库的包设计，如 `net/http`、`database/sql`
- 尝试重构自己的项目，合理划分包边界

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 包 | Package | Go 代码组织的基本单位，一个目录就是一个包 |
| 导出 | Exported | 首字母大写的标识符，包外可以访问 |
| 未导出 | Unexported | 首字母小写的标识符，仅包内可见 |
| 模块 | Module | 一组有版本信息的 Go 包，由 `go.mod` 定义 |
| 导入路径 | Import Path | 导入包时使用的路径，如 `"hello/internal/xxx"` |
| 初始化函数 | Init Function | `init()`，包加载时自动执行 |
| 循环依赖 | Import Cycle | 两个或多个包互相导入，Go 禁止这种结构 |
| 内部包 | Internal Package | `internal/` 目录下的包，外部模块不能导入 |

## 相关资源

- [Go 官方包管理文档](https://go.dev/doc/code#Packages)
- [Go 模块入门教程](https://go.dev/doc/tutorial/create-module)
- [Organizing Go Modules 博客](https://go.dev/blog/module-layout)
- [Go 项目结构标准](https://github.com/golang-standards/project-layout)

[源码](../../internal/basic/packages/packages.go)
