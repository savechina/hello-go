# 函数（Functions）

## 开篇故事

想象你在组装家具。如果所有步骤都写在一张纸上——"拿螺丝、拧木板、装抽屉、贴标签"——你会手忙脚乱。但如果把步骤拆成几个小卡片："组装框架"、"安装抽屉"、"贴标签"，每个卡片只做一件事，整个过程就清晰多了。

Go 的**函数**就是这些小卡片——它们帮你把复杂的程序拆成一个个可理解、可复用、可测试的小单元。

---

## 本章适合谁

如果你想理解如何组织 Go 代码、如何设计函数签名、如何处理错误，本章适合你。你需要理解变量和数据类型，不需要任何函数设计经验。

---

## 你会学到什么

完成本章后，你可以：

1. 定义函数，理解参数（parameters）和返回值（return values）的设计
2. 使用多个返回值，理解 Go 的"结果 + 错误"模式
3. 使用命名返回值（named returns），让返回值语义更清晰
4. 使用可变参数（variadic parameters）处理不定数量的输入
5. 使用闭包（closure）创建有状态的函数

---

## 前置要求

- 理解变量声明（`var` 和 `:=`）
- 理解基础数据类型（int, string, bool）

---

## 第一个例子

让我们从最简单的函数开始：

```go
func greet(name string) string {
    return "Hello, " + name
}

// 调用函数
message := greet("Gopher")
fmt.Println(message)  // 输出：Hello, Gopher
```

**关键概念**：

- `func` - 函数声明关键字
- 参数类型写在参数名**后面**（Go 的特色）
- 返回值类型写在参数列表**后面**

---

## 原理解析

### 1. 函数是组织逻辑的基本单元

Go 鼓励写"小函数"——每个函数只做一件事：

```go
// ❌ 不好：一个函数做太多事
func processUser() {
    // 读取数据库
    // 验证数据
    // 发送邮件
    // 更新日志
    // ... 50 行代码
}

// ✅ 好：拆成小函数
func fetchUser(id int) (User, error) { ... }
func validateUser(u User) error { ... }
func sendWelcomeEmail(u User) error { ... }
func logAction(action string) { ... }
```

**类比**：
> 函数就像厨房里的工具——菜刀切菜、剪刀剪包装、开瓶器开瓶子。每个工具专注一件事，效率最高。

### 2. 多个返回值

Go 允许函数返回多个值，这是它最实用的特性之一：

```go
func divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

// 调用
result, err := divide(10, 3)
```

**"结果 + 错误"模式**是 Go 的标准做法：

| 模式            | 示例                          | 说明                     |
| --------------- | ----------------------------- | ------------------------ |
| 结果 + 错误     | `value, err := parse(s)`      | 最常见                   |
| 结果 + 是否存在 | `value, ok := m[key]`         | map 查找                 |
| 结果 + 是否关闭 | `msg, open := <-ch`           | channel 接收             |
| 多个相关结果    | `quotient, remainder := div()` | 数学运算                 |

### 3. 命名返回值（Named Returns）

你可以给返回值起名字，让它们自带文档：

```go
func rectangleMetrics(width, height float64) (area float64, perimeter float64) {
    area = width * height
    perimeter = 2 * (width + height)
    return  // 裸返回（bare return），自动返回命名变量
}
```

**何时使用命名返回值**：

- ✅ 返回值含义不直观时（如 `area`, `perimeter`）
- ✅ 函数较长，需要文档说明返回值时
- ❌ 简单函数不需要（如 `func add(a, b int) int`）

### 4. 可变参数（Variadic Parameters）

`...T` 让函数接收任意数量的参数：

```go
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

sum(1, 2, 3)       // 6
sum(10, 20)        // 30
sum()              // 0
```

**类比**：
> 可变参数就像一个"无限容量的篮子"——你可以放 0 个、1 个、或任意多个苹果。

### 5. 闭包（Closures）

闭包是"记住外部变量"的函数：

```go
func makeCounter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

counter := makeCounter()
fmt.Println(counter())  // 1
fmt.Println(counter())  // 2
fmt.Println(counter())  // 3
```

**闭包的关键**：内部函数捕获了外部的 `count` 变量，每次调用都会修改它。

**类比**：
> 闭包就像一个带记忆的小盒子——你给它一个初始状态，它每次被调用时都能记住上次做了什么。

---

## 常见错误

### 错误 1: 忽略错误返回值

```go
result, _ := divide(10, 0)  // ❌ 忽略了错误
fmt.Println(result)         // 输出 0，但不知道为什么
```

**症状**：
- 程序行为异常，但找不到原因
- 零值被当作有效值使用

**修复方法**：

总是检查错误：
```go
result, err := divide(10, 0)
if err != nil {
    fmt.Printf("错误: %v\n", err)
    return
}
fmt.Println(result)
```

---

### 错误 2: 闭包捕获循环变量

```go
func main() {
    var funcs []func()
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() {
            fmt.Println(i)  // ❌ 所有函数都打印 3
        }())
    }
}
```

**为什么会这样？**

所有闭包捕获的是**同一个** `i` 变量。循环结束时 `i = 3`，所以所有函数都打印 3。

**修复方法**：

把循环变量作为参数传入：
```go
for i := 0; i < 3; i++ {
    funcs = append(funcs, func(n int) func() {
        return func() {
            fmt.Println(n)  // ✅ 每个函数有自己的 n
        }
    }(i)())
}
```

或者在循环内创建新变量：
```go
for i := 0; i < 3; i++ {
    n := i  // 创建新变量
    funcs = append(funcs, func() {
        fmt.Println(n)  // ✅ 每个函数有自己的 n
    })
}
```

---

### 错误 3: 裸返回导致混淆

```go
func confusing() (result int) {
    if true {
        result = 10
        return  // ✅ 裸返回，返回 10
    }
    return 20  // ❌ 显式返回 20，覆盖了命名返回值
}
```

**修复方法**：

要么全用裸返回，要么全用显式返回，不要混用：
```go
func clear() (result int) {
    if true {
        result = 10
        return  // ✅ 一致
    }
    result = 20
    return  // ✅ 一致
}
```

---

## 动手练习

### 练习 1: 预测输出

不运行代码，预测下面代码的输出：

```go
func swap(a, b string) (string, string) {
    return b, a
}

x, y := swap("hello", "world")
fmt.Println(x, y)
```

<details>
<summary>点击查看答案</summary>

**输出**:
```
world hello
```

**解析**：
1. `swap` 接收两个字符串，返回两个字符串
2. 返回值顺序是 `(b, a)`，所以 `x = "world"`, `y = "hello"`

</details>

---

### 练习 2: 修复错误

下面的代码忽略了错误，请修复：

```go
func parseAge(s string) (int, error) {
    age, err := strconv.Atoi(s)
    return age, err
}

func main() {
    age, _ := parseAge("abc")  // ❌ 忽略错误
    fmt.Printf("年龄: %d\n", age)  // 输出 0
}
```

<details>
<summary>点击查看修复方法</summary>

**修复**：
```go
func main() {
    age, err := parseAge("abc")
    if err != nil {
        fmt.Printf("解析失败: %v\n", err)
        return
    }
    fmt.Printf("年龄: %d\n", age)
}
```

**输出**:
```
解析失败: strconv.Atoi: parsing "abc": invalid syntax
```

</details>

---

### 练习 3: 实现闭包

实现一个 `makeMultiplier` 函数，返回一个闭包，每次调用时将输入乘以固定的倍数：

```go
func makeMultiplier(factor int) func(int) int {
    // 你的代码
}

double := makeMultiplier(2)
triple := makeMultiplier(3)

fmt.Println(double(5))  // 应该输出 10
fmt.Println(triple(5))  // 应该输出 15
```

<details>
<summary>点击查看参考实现</summary>

```go
func makeMultiplier(factor int) func(int) int {
    return func(n int) int {
        return n * factor
    }
}
```

**解析**：
- `factor` 被闭包捕获
- 每次调用返回的函数时，`factor` 保持不变
- 不同的 `makeMultiplier` 调用产生不同的 `factor`

</details>

---

## 故障排查 (FAQ)

### Q: 函数应该写多长？

**A**: 没有硬性限制，但遵循这个原则：

- **理想长度**：10-30 行
- **如果超过 50 行**：考虑拆分
- **如果低于 5 行**：可能太碎了，考虑合并

**判断标准**：如果函数名需要用"和"连接（如 `fetchAndValidateAndSend`），说明它做了太多事。

---

### Q: 什么时候用命名返回值，什么时候不用？

**A**: 遵循这个指南：

- **用命名返回值**：返回值含义不直观、函数较长、需要文档说明
- **不用命名返回值**：简单函数、返回值一目了然

```go
// ✅ 好：命名返回值让含义清晰
func divMod(a, b int) (quotient int, remainder int) { ... }

// ✅ 好：简单函数不需要
func add(a, b int) int { return a + b }

// ❌ 不好：简单函数加命名返回值反而噪音
func add(a, b int) (sum int) { sum = a + b; return }
```

---

### Q: Go 有函数重载（overloading）吗？

**A**: **没有**。Go 不支持函数重载。

```go
// ❌ Go 不允许
func print(s string) { ... }
func print(i int) { ... }
```

**替代方案**：
- 用不同的函数名：`printString(s)`, `printInt(i)`
- 用可变参数：`print(args ...interface{})`
- 用接口：`print(v fmt.Stringer)`

---

## 知识扩展 (选学)

### defer：延迟执行

`defer` 让函数在**当前函数返回前**执行：

```go
func readFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // 在函数返回前关闭文件

    // ... 读取文件
    return nil
}
```

**defer 的三个用途**：
1. 资源清理（关闭文件、数据库连接）
2. 解锁（`defer mu.Unlock()`）
3. 捕获 panic（`defer recover()`）

**执行顺序**：多个 `defer` 按**后进先出**（LIFO）顺序执行：

```go
defer fmt.Println("1")
defer fmt.Println("2")
defer fmt.Println("3")
// 输出：3, 2, 1
```

---

### 函数是一等公民

在 Go 中，函数可以：
- 赋值给变量
- 作为参数传递
- 作为返回值

```go
// 函数作为参数
func apply(f func(int) int, x int) int {
    return f(x)
}

result := apply(func(n int) int { return n * 2 }, 5)  // 10
```

这是中间件（middleware）和装饰器（decorator）模式的基础。

---

## 工业界应用：HTTP 中间件

**场景**：给 HTTP 处理函数添加日志和认证

```go
// 中间件：记录请求日志
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("[%s] %s\n", r.Method, r.URL.Path)
        next(w, r)
    }
}

// 中间件：检查认证
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token != "secret" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

func main() {
    handler := loggingMiddleware(authMiddleware(handleRequest))
    http.ListenAndServe(":8080", handler)
}
```

**为什么这样设计**：
- 每个中间件是一个函数，可复用
- 中间件链可以任意组合
- 核心业务逻辑不受影响

---

## 小结

**核心要点**：

1. **函数应该小且专注** - 每个函数只做一件事
2. **多返回值是 Go 的特色** - 自然表达"结果 + 错误"
3. **总是检查错误** - 不要用 `_` 忽略错误
4. **可变参数处理不定输入** - `...T` 让函数更灵活
5. **闭包捕获外部变量** - 适合计数器和工厂函数

**关键术语**：

- **Parameter (参数)**: 函数接收的输入
- **Return Value (返回值)**: 函数的输出
- **Named Return (命名返回值)**: 有名字的返回值，自带文档
- **Variadic (可变参数)**: 接收任意数量参数的函数
- **Closure (闭包)**: 记住外部变量状态的函数
- **Defer**: 延迟执行，常用于资源清理

**下一步**：

- 继续：[流程控制](flowcontrol.md)
- 回顾：[阶段复习](review-basic.md)

---

## 术语表

| English          | 中文       |
| ---------------- | ---------- |
| Function         | 函数       |
| Parameter        | 参数       |
| Return Value     | 返回值     |
| Named Return     | 命名返回值 |
| Variadic         | 可变参数   |
| Closure          | 闭包       |
| Defer            | 延迟执行   |
| Middleware       | 中间件     |

---

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/functions/functions.go)
