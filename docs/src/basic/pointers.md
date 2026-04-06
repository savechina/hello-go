# 指针（Pointers）

## 开篇故事

想象你有两个笔记本，一个是自己的，一个是朋友的。朋友说："帮我记个电话号码。"你有两种选择：

1. **复制一份**：把朋友的本子拿过来，抄下所有内容到自己本子上，然后改。改完后，朋友的本子还是原来的内容——白忙活了。
2. **直接修改**：让朋友把本子递过来，你在上面直接写。改完后，朋友拿回去就能看到新内容。

指针就是第二种方式。`&` 是"给我你的本子"（取地址），`*` 是"打开本子写字"（解引用）。如果不用指针，函数参数永远是"复制一份"，函数内部修改不影响外部。用了指针，函数就能"拿到你的本子"，直接修改同一份数据。

初学者觉得指针抽象，是因为它涉及"内存地址"这个概念。但换个角度想：指针就是个"遥控器"，按按钮（解引用）就能控制电视（原变量）。你不需要知道电视内部电路怎么工作，只需要知道遥控器能干什么。

本章用最实用的方式讲解指针：什么时候用、怎么用、怎么避免踩坑。

## 本章适合谁

- 你见过 `*int`、`&value` 这种语法，但不清楚它们到底干嘛的
- 你写过函数，发现修改参数不影响外部变量
- 你想理解方法接收者 `func (w *wallet) Deposit()` 为什么用 `*`
- 你遇到过 `nil pointer dereference` 错误，想学会避免它

如果你还没学过 Go 变量和函数，建议先看 [变量](./variables.md)、[函数](./functions.md)；如果你想深入理解内存模型，可以学习 [内存管理](../advanced/memory.md)。

## 你会学到什么

学完本章，你将能够：

1. 理解指针的本质：存储变量内存地址的特殊变量
2. 使用 `&` 取地址、`*` 解引用，在函数间传递指针
3. 理解指针接收者（pointer receiver）的作用和语法
4. 安全处理 nil 指针，避免运行时 panic
5. 判断何时应该用指针、何时用值传递

## 前置要求

在开始之前，你需要：

- **理解变量**：知道变量存储数据，有类型和值
- **理解函数参数**：知道参数是"传值"的，函数内部修改不影响外部
- **理解结构体**：指针经常和结构体一起使用，特别是方法接收者
- **基础语法**：`if`、`return`、`fmt.Println` 等基本概念

如果这些概念还不熟悉，建议先阅读：[变量与常量](./variables.md)、[结构体](./structs.md)。

## 第一个例子

让我们从一个最简单的例子开始：修改一个变量的值。

### 不用指针会怎么样

```go
func tryModify(value int) {
    value = 100  // 只是修改了副本
}

func main() {
    x := 10
    tryModify(x)
    fmt.Println(x)  // 输出：10，原值没变
}
```

函数参数 `value` 是 `x` 的副本，改了也白改。

### 使用指针

```go
func modifyWithPointer(pointer *int) {
    *pointer = 100  // 通过指针修改原值
}

func main() {
    x := 10
    modifyWithPointer(&x)  // 传入 x 的地址
    fmt.Println(x)  // 输出：100，原值被修改
}
```

**关键步骤**：
1. `&x`：取 `x` 的地址，类型是 `*int`
2. `pointer *int`：函数参数声明为指针类型
3. `*pointer = 100`：解引用，修改地址指向的值

这个例子展示了指针的核心价值：**让函数能够修改调用方的变量**。

## 原理解析

### 1. 地址和解引用

每个变量在内存中都有一个地址。`&` 运算符可以获取这个地址：

```go
value := 10
pointer := &value  // pointer 存储 value 的内存地址
```

`pointer` 是一个指针变量，它的类型是 `*int`（指向 int 的指针）。

要访问或修改地址中的值，需要用 `*` 解引用：

```go
*pointer = 15      // 修改原值
fmt.Println(value) // 输出：15
fmt.Println(*pointer) // 输出：15，和 value 一样
```

**关键理解**：
- `pointer` 的值是"地址"（比如 `0xc000016080`）
- `*pointer` 的值是"地址里存储的数据"（比如 `15`）

### 2. 指针接收者（Pointer Receiver）

方法可以用指针作为接收者，这样方法就能修改对象状态：

```go
type wallet struct {
    balance int
}

func (w *wallet) Deposit(amount int) {
    if w == nil {
        return
    }
    w.balance += amount
}

func (w *wallet) Balance() int {
    if w == nil {
        return 0
    }
    return w.balance
}
```

调用时：

```go
account := &wallet{}  // 创建指针
account.Deposit(30)
account.Deposit(12)
fmt.Println(account.Balance())  // 输出：42
```

**为什么要用指针接收者**：
- **修改状态**：值接收者（`func (w wallet)`）修改的是副本
- **避免复制**：大结构体用指针接收者更高效
- **一致性**：如果一个方法用指针接收者，所有方法都应该用

### 3. nil 指针和安全检查

指针可以是 `nil`，表示"不指向任何东西"：

```go
var nobody *learner  // nobody 是 nil
var broken *wallet   // broken 是 nil
```

直接解引用 nil 指针会 panic：

```go
fmt.Println(*nobody)  // ✗ panic: invalid memory address or nil pointer dereference
```

**安全做法**：先检查是否为 nil

```go
func safeLearnerName(item *learner) string {
    if item == nil {
        return "nil learner"
    }
    return item.name
}

func (w *wallet) Balance() int {
    if w == nil {
        return 0  // 返回默认值，而不是 panic
    }
    return w.balance
}
```

这种模式很常见：**方法对 nil 接收者有良好行为**。

### 4. 指针作为函数参数

函数参数用指针，可以修改多个变量或避免大对象复制：

```go
// 交换两个变量的值
func swapValues(left *int, right *int) bool {
    if left == nil || right == nil {
        return false
    }
    *left, *right = *right, *left
    return true
}

// 修改字符串
func renameWithPointer(target *string, next string) bool {
    if target == nil {
        return false
    }
    *target = next
    return true
}
```

调用：

```go
a := 10
b := 20
swapValues(&a, &b)
fmt.Printf("a=%d, b=%d\n", a, b)  // 输出：a=20, b=10

name := "Alice"
renameWithPointer(&name, "Bob")
fmt.Println(name)  // 输出：Bob
```

### 5. 指针的零值

指针的零值是 `nil`：

```go
var p *int
fmt.Println(p == nil)  // true
```

创建指针有三种方式：

```go
// 方式 1：用 & 取地址
value := 100
p1 := &value

// 方式 2：用 new()
p2 := new(int)  // *p2 是 int 的零值 0
*p2 = 200

// 方式 3：结构体直接用复合字面量
w := &wallet{balance: 100}
```

## 常见错误

### 错误 1：忘记解引用

```go
func wrong(p *int) {
    p = 100  // ✗ 类型不匹配：不能把 int 赋给 *int
}

func right(p *int) {
    *p = 100  // ✓ 解引用后赋值
}
```

**修复**：用 `*p` 而不是 `p`。

### 错误 2：忽略 nil 检查

```go
type user struct {
    name string
}

func getName(u *user) string {
    return u.name  // ✗ 如果 u 是 nil，会 panic
}

// 修复
func getName(u *user) string {
    if u == nil {
        return ""
    }
    return u.name
}
```

**最佳实践**：导出函数对 nil 输入应该有良好行为。

### 错误 3：不必要的指针

```go
// 过度使用指针
func add(a *int, b *int) *int {
    result := *a + *b
    return &result  // 返回局部变量地址（虽然 Go 有逃逸分析，但不推荐）
}

// 更简洁的写法
func add(a int, b int) int {
    return a + b
}
```

**原则**：不需要修改参数时，用值传递。

## 动手练习

### 练习 1：计数器

实现一个计数器类型，有 `Increment()`、`Decrement()`、`Value()` 方法，要求用指针接收者：

```go
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    // 你的代码
}

func main() {
    c := &Counter{}
    c.Increment()
    c.Increment()
    c.Decrement()
    fmt.Println(c.Value())  // 应该输出：1
}
```

<details>
<summary>参考答案</summary>

```go
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    if c == nil {
        return
    }
    c.value++
}

func (c *Counter) Decrement() {
    if c == nil {
        return
    }
    c.value--
}

func (c *Counter) Value() int {
    if c == nil {
        return 0
    }
    return c.value
}
```

</details>

### 练习 2：指针交换器

写一个函数，交换两个字符串指针指向的内容：

```go
func swapStrings(a *string, b *string) {
    // 你的代码
}

func main() {
    x := "hello"
    y := "world"
    swapStrings(&x, &y)
    fmt.Println(x, y)  // 应该输出：world hello
}
```

<details>
<summary>参考答案</summary>

```go
func swapStrings(a *string, b *string) {
    if a == nil || b == nil {
        return
    }
    *a, *b = *b, *a
}
```

</details>

### 练习 3：安全访问嵌套指针

有一个结构体包含指针字段，写一个安全函数访问深层嵌套的值：

```go
type Address struct {
    City string
}

type Person struct {
    Name    string
    Address *Address
}

func getCity(p *Person) string {
    // 你的代码：要处理 p 为 nil、p.Address 为 nil 的情况
}
```

<details>
<summary>参考答案</summary>

```go
func getCity(p *Person) string {
    if p == nil {
        return ""
    }
    if p.Address == nil {
        return ""
    }
    return p.Address.City
}

// 或者用一行
func getCity(p *Person) string {
    if p != nil && p.Address != nil {
        return p.Address.City
    }
    return ""
}
```

</details>

## 故障排查 (FAQ)

### Q1: 什么时候应该用指针，什么时候用值？

**A**: 遵循这些原则：

**用指针的情况**：
- 需要修改参数或接收者
- 结构体很大（比如超过 3 个字段），复制成本高
- 需要表示"不存在"（nil）
- 方法需要保持一致性（如果一个用指针，全部用指针）

**用值的情况**：
- 基本类型（int、string、bool）
- 小结构体（1-2 个字段）
- 不需要修改，也不想让调用方看到变化
- 类型本身是引用类型（map、slice、channel）

**经验法则**：如果不确定，先看标准库同类型怎么处理。

### Q2: `nil` 指针一定有问题吗？

**A**: 不一定。Go 的风格鼓励**对 nil 友好**：

```go
// 好的设计
func (w *wallet) Balance() int {
    if w == nil {
        return 0  // 返回合理的零值
    }
    return w.balance
}

// 调用方不需要担心
var w *wallet
fmt.Println(w.Balance())  // 输出：0，不会 panic
```

**坏的设计**是让调用方必须检查 nil，否则就 panic。

### Q3: 指针和内存泄漏有关系吗？

**A**: Go 有垃圾回收（GC），不用担心忘记释放指针。但要注意：

```go
// 可能的问题：意外保持引用
type Cache struct {
    data map[string]*LargeObject
}

// 删除时只删了 map 里的引用，但其他地方可能还持有指针
delete(c.data, "key")
```

**建议**：
- 不要过度使用指针，特别是短生命周期的对象
- 注意循环引用（A 指向 B、B 指向 A），GC 能处理但可能影响性能
- 用 `go tool pprof` 检测内存问题

## 知识扩展 (选学)

### 1. 指针的指针

指针本身也是变量，也可以取地址：

```go
x := 10
p := &x      // *int
pp := &p     // **int

fmt.Println(**pp)  // 输出：10
```

这种场景很少见，通常用在需要修改指针本身的情况。

### 2. 方法值和方法表达式

Go 有高级特性，可以把方法绑定到变量：

```go
w := &wallet{balance: 100}

// 方法值（method value）
deposit := w.Deposit
deposit(50)  // 等价于 w.Deposit(50)

// 方法表达式（method expression）
wallet.Deposit(w, 50)  // 接收者作为第一个参数
```

这在函数式编程或回调中很有用。

### 3. 逃逸分析（Escape Analysis）

Go 编译器会决定变量分配在栈上还是堆上：

```go
func localPointer() *int {
    x := 10
    return &x  // x 会"逃逸"到堆上，不会变成悬垂指针
}
```

用 `go build -gcflags="-m"` 可以看到分析结果。

### 4. 指针和接口

接口内部存储指针时，nil 检查要小心：

```go
type Speaker interface {
    Speak()
}

type Dog struct{}

func (d *Dog) Speak() { fmt.Println("woof") }

var s Speaker = (*Dog)(nil)  // 接口值是 (*Dog, nil)
fmt.Println(s == nil)  // false！接口本身不是 nil
```

**规则**：只有接口值和动态值都 nil 时，`interface == nil` 才为 true。

### 5. `unsafe.Pointer`（危险操作）

`unsafe` 包允许绕过类型系统：

```go
import "unsafe"

x := 10
p := unsafe.Pointer(&x)  // 可以转换为任何指针类型
```

**警告**：这会破坏类型安全，只在特殊场景使用（如系统编程、序列化）。

## 工业界应用

### 场景：数据库连接池

在 Web 服务中，数据库连接通常是共享资源，用指针传递：

```go
type Database struct {
    conn *sql.DB
}

func NewDatabase(dsn string) (*Database, error) {
    conn, err := sql.Open("mysql", dsn)
    if err != nil {
        return nil, err
    }
    return &Database{conn: conn}, nil
}

func (db *Database) Query(ctx context.Context, query string) (*Rows, error) {
    if db == nil || db.conn == nil {
        return nil, errors.New("database not initialized")
    }
    return db.conn.QueryContext(ctx, query)
}
```

**为什么用指针**：
- 连接是共享资源，不能被复制
- 需要表示"未初始化"状态（nil）
- 避免每次查询都复制大的连接对象

### 场景：配置对象

配置通常在启动时加载，运行中可能被热更新：

```go
type Config struct {
    Port     int
    LogLevel string
    mu       sync.RWMutex
}

func (c *Config) GetLogLevel() string {
    if c == nil {
        return "info"  // 默认值
    }
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.LogLevel
}

func (c *Config) SetLogLevel(level string) {
    if c == nil {
        return
    }
    c.mu.Lock()
    defer c.mu.Unlock()
    c.LogLevel = level
}
```

**指针的作用**：
- 所有模块共享同一份配置
- 读写需要加锁，指针确保锁的是同一个对象
- nil 检查提供安全的降级行为

### 真实案例：标准库 `bytes.Buffer`

看看标准库怎么用指针：

```go
type Buffer struct {
    buf       []byte
    off       int
    lastRead  readOp
}

func NewBuffer() *Buffer {
    return &Buffer{}
}

func (b *Buffer) Write(p []byte) (int, error) {
    // 修改内部 buf
}

func (b *Buffer) Bytes() []byte {
    return b.buf[b.off:]
}
```

`NewBuffer` 返回指针，因为：
- Buffer 内部有切片，复制没有意义
- Write 方法需要修改内部状态
- 避免每次调用都分配新 Buffer

### 真实案例：HTTP 处理器

```go
type Handler struct {
    db *Database
}

func NewHandler(db *Database) *Handler {
    return &Handler{db: db}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if h == nil || h.db == nil {
        http.Error(w, "service unavailable", 503)
        return
    }
    // 处理请求
}
```

依赖注入（Dependency Injection）模式的核心就是指针对象的传递。

## 小结

本章我们学习了：

1. **地址和解引用**：`&` 取地址，`*` 访问地址中的值
2. **指针接收者**：方法用 `*T` 接收者可以修改对象状态
3. **nil 指针**：指针可以是 nil，使用前要检查
4. **指针参数**：函数参数用指针可以修改调用方变量
5. **使用场景**：修改、共享、避免复制、表示"不存在"

关键术语：
- **指针（Pointer）**：存储内存地址的变量，类型如 `*int`
- **取地址（Address-of）**：`&` 运算，获取变量地址
- **解引用（Dereference）**：`*` 运算，访问地址中的值
- **指针接收者（Pointer Receiver）**：方法接收者是指针类型
- **nil**：指针的零值，表示"不指向任何东西"

下一步建议：
- 阅读 Go 官方文档 "Effective Go" 的指针部分
- 学习 [接口](./interfaces.md)，理解接口和指针的关系
- 用 `go vet` 检查代码中的指针问题

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 指针 | Pointer | 存储变量内存地址的特殊变量 |
| 取地址 | Address-of | 使用 `&` 获取变量的内存地址 |
| 解引用 | Dereference | 使用 `*` 访问指针指向的值 |
| 指针接收者 | Pointer Receiver | 方法接收者声明为指针类型，如 `(t *T)` |
| nil | nil | Go 的零值，指针的默认值是 nil |
| 值传递 | Pass by Value | 函数参数是副本，修改不影响原值 |
| 引用传递 | Pass by Reference | 通过指针让函数修改原值 |
| 逃逸分析 | Escape Analysis | 编译器决定变量分配在栈上还是堆上 |

## 相关资源

- [Effective Go - Pointers](https://go.dev/doc/effective_go#pointers)
- [Go 官方博客：指针入门](https://blog.golang.org/pointers)
- [Understanding Go Pointers 教程](https://medium.com/@meeusdylan/understanding-pointers-in-go-6ebc9b9d7a91)
- [Go 指针最佳实践](https://github.com/golang/go/wiki/CodeReviewComments#pointers)

[源码](../../internal/basic/pointers/pointers.go)
