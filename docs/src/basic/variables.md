# 变量与表达式（Variables & Expressions）

## 开篇故事

想象你有一个工具箱，里面装着各种工具：螺丝刀、锤子、尺子。你给每个工具贴上标签，下一次需要时就知道去哪里找。Go 中的**变量**就像这些贴标签的工具箱——它们帮你存储和管理程序中的数据。**常量**则是那些你钉在墙上的工具——一旦放好，就不会再移动。

---

## 本章适合谁

如果你是 Go 初学者，想理解如何存储数据、声明常量和进行基本计算，本章适合你。这是所有编程的基础，即使你是第一次接触编程也能理解。

---

## 你会学到什么

完成本章后，你可以：

1. 使用 `var` 关键字声明变量，理解何时需要显式写类型
2. 使用 `:=` 短变量声明，理解类型推断（type inference）
3. 使用 `const` 声明常量，理解不可变值的意义
4. 使用 `iota` 生成连续的常量编号
5. 区分"应该用变量"和"应该用常量"的场景

---

## 前置要求

本章是 Go 的第一章，不需要前置知识。如果你有任意编程基础（Python、JavaScript、Java 等）会更容易理解。

---

## 第一个例子

让我们从最简单的变量声明开始：

```go
var language string = "Go"
var lessonCount int = 12
var ready bool = true
```

**关键概念**：

- `var` - 声明变量的关键字
- 类型写在变量名**后面**（这是 Go 和其他语言的重要区别）
- 每个变量声明后都有初始值

---

## 原理解析

### 1. 变量声明（var）

在 Go 中，`var` 是最基础的变量声明方式：

```go
var language string = "Go"
var lessonCount int = 12
```

**为什么 Go 要把类型写在后面？**

- 其他语言（C/Java）：`string language = "Go";`
- Go：`var language string = "Go"`

Go 的设计者认为，当类型显而易见时，你根本不需要写它。这让 `var` 和 `:=` 的视觉模式更一致。

**类比**：
> 就像你填写表格——先写名字，再写类型（"姓名：张三"），而不是先写类型再写名字。

### 2. 短变量声明（:=）

当类型显而易见时，Go 允许你省略 `var` 和类型：

```go
name := "Alice"
age := 30
```

**`var` vs `:=` 的选择指南**：

| 场景              | 推荐写法 | 原因                       |
| ----------------- | -------- | -------------------------- |
| 函数内部，类型明显  | `:=`     | 简洁，最常见               |
| 想突出类型          | `var`    | 让读者注意到类型           |
| 包级别（函数外部）  | `var`    | `:=` 只能在函数内部使用    |
| 需要零值（zero value） | `var` | `var x int` 会得到 `0`     |

### 3. 类型推断（Type Inference）

Go 会根据右侧的值自动推断变量类型：

```go
total := 3          // int
progress := 75.5    // float64
note := "go"        // string
ready := true       // bool
```

**重要**：虽然 Go 帮你推断类型，但你依然需要知道最终推断出的是什么类型。因为类型会影响运算、函数调用和接口匹配。

### 4. 常量（Constants）

常量是**永远不变**的值：

```go
const courseName = "hello-go"
const maxRetries = 3
```

**常量 vs 变量**：

| 特征           | 变量 (`var`/`:=`) | 常量 (`const`)        |
| -------------- | ----------------- | --------------------- |
| 可修改         | ✅ 是（除非用 const） | ❌ 否                 |
| 运行时确定     | ✅ 是               | ❌ 否（编译期已知）     |
| 可以使用函数值 | ✅ 是               | ❌ 否（只能用字面量）   |
| 生命周期       | 作用域内           | 整个程序运行期间       |

**何时使用常量**：

- 配置值（最大重试次数、超时时间）
- 状态名（"draft", "review", "published"）
- 数学常数

### 5. iota 生成连续常量

`iota` 是 Go 用来生成连续常量值的内建标识符：

```go
const (
    stageDraft = iota    // 0
    stageReview          // 1
    stagePublished       // 2
)
```

**类比**：
> 就像自动编号的发票——你不需要手动写 1、2、3，机器帮你递增。

**为什么用 iota 而不是手写数字？**

- 不容易出错（不会漏掉某个编号）
- 更容易维护（中间插入新状态时，后面的自动调整）
- 意图更清晰（读者一眼就知道这是连续编号）

---

## 常见错误

### 错误 1: 在函数外部使用 :=

```go
package main

x := 5  // ❌ 编译错误！

func main() {}
```

**编译器输出**:
```
syntax error: non-declaration statement outside function body
```

**修复方法**：

在包级别使用 `var`：
```go
package main

var x = 5  // ✅

func main() {}
```

---

### 错误 2: 常量使用运行时值

```go
const now = time.Now()  // ❌ 编译错误！
```

**编译器输出**:
```
const initializer time.Now() is not a constant
```

**修复方法**：

改用 `var`：
```go
var now = time.Now()  // ✅
```

---

### 错误 3: 未使用变量的警告

```go
func main() {
    unused := 5  // ⚠️ 编译错误！Go 不允许未使用的变量
}
```

**编译器输出**:
```
unused declared and not used
```

**修复方法**：

使用前缀下划线或真正使用它：
```go
func main() {
    _ = 5  // ✅ 编译器知道你是故意的
}
```

---

## 动手练习

### 练习 1: 预测输出

不运行代码，预测下面代码的输出：

```go
x := 5
x = x + 1
{
    x := x * 2
    fmt.Println("内部：", x)
}
fmt.Println("外部：", x)
```

<details>
<summary>点击查看答案</summary>

**输出**:
```
内部： 12
外部： 6
```

**解析**：
1. `x = 5` - 初始值
2. `x = 5 + 1 = 6` - 修改 x
3. 内部作用域：`x := 6 * 2 = 12` - 新变量遮蔽了外部 x
4. 内部作用域结束，内部 x 失效
5. 外部 x 仍然是 6

</details>

---

### 练习 2: 修复错误

下面的代码有编译错误，请修复：

```go
const maxUsers = 100
maxUsers = 200  // ❌ 错误
```

<details>
<summary>点击查看修复方法</summary>

**修复**：
```go
var maxUsers = 100  // 改为 var
maxUsers = 200      // ✅ 现在可以修改了
```

**或者**，如果你确实需要常量，就不要修改它：
```go
const maxUsers = 100
// maxUsers = 200  // ❌ 常量不能修改
newMax := 200       // ✅ 创建新变量
```

</details>

---

### 练习 3: 使用 iota

改写下面的代码，使用 iota 代替手动编号：

```go
const (
    statusPending = 0
    statusActive = 1
    statusCompleted = 2
    statusArchived = 3
)
```

<details>
<summary>点击查看参考实现</summary>

```go
const (
    statusPending = iota  // 0
    statusActive          // 1
    statusCompleted       // 2
    statusArchived        // 3
)
```

**好处**：
- 不需要手动写数字
- 中间插入新状态时，后面的自动调整
- 意图更清晰

</details>

---

## 故障排查 (FAQ)

### Q: 什么时候应该用 `var`，什么时候用 `:=`？

**A**: 遵循这个原则：

- **函数内部，类型明显** → 用 `:=`（90% 的情况）
- **想突出类型** → 用 `var`
- **包级别（函数外部）** → 只能用 `var`
- **需要零值语义** → 用 `var`（如 `var count int` 得到 `0`）

示例：
```go
// ✅ 好的实践
var config *Config  // 包级别，突出类型

func main() {
    name := "hello"  // 类型明显是 string
    var count int    // 需要零值 0
}
```

---

### Q: 为什么 Go 不允许未使用的变量？

**A**: 这是 Go 的设计哲学——未使用的变量通常是 bug 的信号。

- **C/Java/Python**：未使用变量只是警告
- **Go**：未使用变量是编译错误

**好处**：
1. 减少代码噪音（没有"死代码"）
2. 避免拼写错误（`userName` vs `userNmae`）
3. 强制你清理不需要的代码

---

### Q: `const` 和 `var` 的性能有区别吗？

**A**: 有，但通常可以忽略。

- `const` 在编译期求值，零运行时开销
- `var` 在运行期初始化

**实际影响**：对于简单类型（int, string），差异在纳秒级别，不需要担心。

---

## 知识扩展 (选学)

### 零值（Zero Value）

Go 中每个类型都有一个"零值"——当你声明变量但不赋值时的默认值：

```go
var i int      // 0
var f float64  // 0.0
var s string   // ""
var b bool     // false
var p *int     // nil
```

**为什么 Go 要设计零值？**

- 避免未初始化变量的 bug（其他语言中常见）
- 简化代码（不需要处处检查 null）
- 让 `var` 声明更简洁

---

### 变量遮蔽（Shadowing）

Go 允许在内部作用域用 `:=` 创建同名变量——新变量会"遮蔽"旧变量：

```go
x := 5
{
    x := 10  // 新 x 遮蔽了外部 x
    fmt.Println(x)  // 10
}
fmt.Println(x)  // 5
```

**遮蔽的优势**：
- 可以改变类型
- 可以复用名称（代码更简洁）
- 在不同作用域有不同含义

**遮蔽的风险**：
- 如果遮蔽让代码更难理解，使用不同的名称

---

## 工业界应用：配置管理

**场景**：Web 服务器配置

```go
const (
    defaultPort = 8080
    defaultHost = "127.0.0.1"
    maxConnections = 1000
)

func main() {
    // 配置在初始化后不应该改变
    port := defaultPort
    host := defaultHost

    fmt.Printf("服务器启动在 %s:%d\n", host, port)
}
```

**为什么常量很重要**：
- 防止运行中意外修改配置
- 集中定义，易于修改
- 编译器保证配置不会被篡改

---

## 小结

**核心要点**：

1. **`var` 是最基础的声明方式** - 可以在任何地方使用
2. **`:=` 是最常见的写法** - 只能在函数内部，依赖类型推断
3. **`const` 表达不变的值** - 编译期已知，运行时不可修改
4. **`iota` 生成连续常量** - 适合状态编号、枚举风格常量
5. **Go 不允许未使用的变量** - 这是编译错误，不是警告

**关键术语**：

- **Type Inference (类型推断)**: 编译器根据右侧值自动推断变量类型
- **Zero Value (零值)**: 变量声明但未赋值时的默认值
- **Shadowing (遮蔽)**: 在内部作用域用同名变量覆盖外部变量
- **iota**: Go 内建的连续常量生成器

**下一步**：

- 继续：[基础数据类型](datatype.md)
- 回顾：[阶段复习](review-basic.md)

---

## 术语表

| English            | 中文     |
| ------------------ | -------- |
| Variable           | 变量     |
| Constant           | 常量     |
| Type Inference     | 类型推断 |
| Zero Value         | 零值     |
| Short Declaration  | 短变量声明 |
| Shadowing          | 遮蔽     |

---

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/variables/variables.go)
