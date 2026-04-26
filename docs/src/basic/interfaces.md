# 接口（Interfaces）

## 开篇故事

想象你去一家多功能餐厅吃饭。你不需要知道厨师是用燃气灶还是电磁炉做菜，不需要知道服务员是用纸笔还是平板点餐，不需要知道收银员是用计算器还是电脑结账。你只需要知道：**厨师会做菜**、**服务员会点餐**、**收银员会结账**。这些"会做什么"就是接口。

在编程中，接口（Interface）定义的是**行为（behavior）**而不是**实现（implementation）**。Go 的接口设计尤其独特：它不需要你显式声明"我实现了这个接口"，只要你的类型有接口要求的所有方法，就**自动实现**了该接口。这种"隐式实现"让代码更灵活、更易于测试、更符合"面向接口编程"的原则。

接口不是抽象的学术概念，而是每天都在用的工具。当你写 `fmt.Println(x)` 时，`x` 可以是字符串、数字、自定义结构体——因为它们都隐式实现了 `fmt.Stringer` 接口。当你用 `io.Copy(dst, src)` 时，不关心 `dst` 是文件、网络还是内存，只要它实现了 `io.Writer`。理解接口，就是理解 Go 的"多态"思维。

## 本章适合谁

- 已经会写结构体，想学习抽象和复用机制的 Go 初学者
- 从 Java/Python 转来 Go，想理解隐式接口的开发者
- 对 `io.Writer`、`io.Reader` 等标准库接口感到困惑的工程师
- 想提高代码可测试性和模块化水平的程序员

## 你会学到什么

完成本章后，你将能够：

1. **理解隐式实现（implicit implementation）**：解释为什么 Go 不需要implements 关键字，以及带来的灵活性
2. **设计和实现小接口**：遵循"最少方法原则"设计灵活的接口
3. **使用 io.Writer/io.Reader 模式**：将标准库接口应用到自定义类型
4. **运用空接口和类型断言**：安全处理任意类型，理解类型switch
5. **用接口解耦依赖**：编写可测试、可替换的模块化代码

## 前置要求

- 已经掌握结构体（Structs）的定义和方法绑定
- 理解指针的基本概念
- 了解函数作为一等公民的特性
- 知道什么是面向对象编程中的"多态"

## 第一个例子

让我们从一个简单的问候场景开始：

```go
package main

import "fmt"

type Greeter interface {
	greet() string
}

type Robot struct {
	Name string
}

func (r Robot) greet() string {
	return fmt.Sprintf("%s says hello", r.Name)
}

func announceGreeter(g Greeter) string {
	return g.greet()
}

func main() {
	r := Robot{Name: "R2"}
	fmt.Println(announceGreeter(r))
	// 输出：R2 says hello
}
```

这个例子展示了接口的核心：**定义行为**（`greet()`）、**隐式实现**（`Robot` 没有写 `implements`）、**面向接口编程**（`announceGreeter` 只关心行为，不关心具体类型）。

## 原理解析

### 1. 隐式实现（Implicit Implementation）

Go 的接口实现不需要关键字：

```go
type Writer interface {
	Write([]byte) (int, error)
}

type FileWriter struct{}

// ✅ 自动实现 Writer，不需要写"implements"
func (f FileWriter) Write(data []byte) (int, error) {
	// ... 写入文件
	return len(data), nil
}
```

**为什么这样设计？**

- **解耦**：实现方不需要知道接口的存在，调用方定义需要什么行为
- **灵活**：同一个类型可以"实现"无数个接口，无需预先声明
- **简洁**：没有冗余的 `implements` 关键字，代码更清爽

**对比 Java**：
```java
// Java 必须显式声明
public class Robot implements Greeter { ... }
```

```go
// Go 自动实现
type Robot struct{}
func (r Robot) greet() string { return "hello" }  // 自动满足 Greeter
```

### 2. 小接口（Small Interface）哲学

Go 鼓励设计**非常小**的接口，通常只有 1-2 个方法：

```go
// io.Writer 只有 1 个方法
type Writer interface {
	Write(p []byte) (n int, err error)
}

// io.Reader 只有 1 个方法
type Reader interface {
	Read(p []byte) (n int, err error)
}

// fmt.Stringer 只有 1 个方法
type Stringer interface {
	String() string
}

// error 只有 1 个方法
type error interface {
	Error() string
}
```

**为什么小接口更好？**

- **更容易实现**：1 个方法比 10 个方法容易实现得多
- **更灵活**：可以组合多个小接口形成大接口
- **更专注**：每个接口只代表一种能力

### 3. io.Writer 模式：标准库的典范

`io.Writer` 是 Go 最重要的接口之一：

```go
func writeLogLine(writer io.Writer, level string, message string) error {
	_, err := fmt.Fprintf(writer, "%s: %s", strings.ToUpper(level), message)
	return err
}

// 可以传入任何实现了 io.Writer 的东西
var buf bytes.Buffer
writeLogLine(&buf, "INFO", "server started")

file, _ := os.Open("log.txt")
writeLogLine(file, "INFO", "server started")

writeLogLine(os.Stdout, "INFO", "server started")
```

**关键洞察**：`writeLogLine` 不关心写入到哪里，只关心**能否写入**。这就是面向接口编程的精髓。

### 4. 空接口（Empty Interface）与类型断言

空接口 `interface{}`（或别名 `any`）不包含任何方法，因此**所有类型都自动实现了它**：

```go
func inspectValue(value any) string {
	switch typed := value.(type) {
	case string:
		return fmt.Sprintf("string => %s", strings.ToUpper(typed))
	case int:
		// 类型断言的短变量形式
		return fmt.Sprintf("int => %d", typed*2)
	case greeter:
		return fmt.Sprintf("greeter => %s", typed.greet())
	default:
		return fmt.Sprintf("unknown => %T", value)
	}
}

func inspectValues(values []any) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, inspectValue(value))
	}
	return strings.Join(parts, " | ")
}

// 使用示例
result := inspectValues([]any{"go", 7, robot{name: "Mika"}})
// 输出：string => GO | int => 14 | greeter => Mika says hello
```

**类型断言（Type Assertion）的两种写法**：

```go
// 方式 1：类型分支（推荐）
switch v := value.(type) {
case int:
	fmt.Println(v * 2)
}

// 方式 2：断言检查
number, ok := value.(int)
if !ok {
	return "assertion failed"
}
```

### 5. 接口组合（Interface Composition）

可以组合多个接口形成新接口：

```go
type ReadWriter interface {
	Reader  // 嵌入 io.Reader
	Writer  // 嵌入 io.Writer
}

// 等价于：
type ReadWriter interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
}
```

**标准库示例**：
```go
type ReadWriter interface {
	Reader
	Writer
}

type ReadWriteCloser interface {
	Reader
	Writer
	Closer
}
```

这种组合机制让接口极其灵活，可以根据需要拼装能力。

## 常见错误

### 错误 1：接口污染（Interface Pollution）

```go
// ❌ 错误：过早定义大接口
type UserServiceInterface interface {
	CreateUser(name string) (*User, error)
	GetUser(id string) (*User, error)
	UpdateUser(id string, name string) error
	DeleteUser(id string) error
	ListUsers() ([]User, error)
	// ... 10 个方法
}

// ✅ 正确：用小接口或直接用结构体
type UserService struct {
	db *Database
}

func (s *UserService) CreateUser(name string) (*User, error) {
	// ...
}
```

**原则**：**不要预先设计接口**，等有两个实现时再提取接口。

### 错误 2：返回值是指针还是接口？

```go
type Config struct {
	Name string
}

// ❌ 错误：返回接口但只有一个实现
func NewConfig() ConfigInterface {
	return &Config{Name: "default"}
}

// ✅ 正确：直接返回具体类型
func NewConfig() *Config {
	return &Config{Name: "default"}
}
```

**经验法则**：只有在需要多态或解耦时才用接口作为返回类型。

### 错误 3：nil 接口值的陷阱

```go
type Printer interface {
	Print()
}

type FilePrinter struct{}

func (f *FilePrinter) Print() {}

// ❌ 陷阱：这会打印"<nil>"而不是"nil"
var p *FilePrinter = nil
var i Printer = p  // i 不是 nil！是 (type=*FilePrinter, value=nil)

if i == nil {
	fmt.Println("is nil")
} else {
	fmt.Println("not nil")  // 会打印这行
}

// ✅ 正确：直接比较具体类型
var p2 *FilePrinter = nil
if p2 == nil {
	fmt.Println("is nil")
}
```

**本质原因**：接口值包含**类型**和**值**两个部分，只有两者都是 nil 时接口才是 nil。

## 动手练习

### 练习 1：预测输出结果

```go
type Formatter interface {
	Format() string
}

type JSON struct{}
func (j JSON) Format() string { return "JSON" }

type XML struct{}
func (x XML) Format() string { return "XML" }

func process(f Formatter) string {
	return f.Format()
}

func main() {
	formatters := []Formatter{JSON{}, XML{}}
	for _, f := range formatters {
		fmt.Print(process(f) + " ")
	}
}
// 问：输出是什么？
```

<details>
<summary>点击查看答案</summary>

```
输出：JSON XML
```

**解析**：`JSON` 和 `XML` 都自动实现了 `Formatter` 接口，可以统一处理。

</details>

### 练习 2：修复错误代码

下面的代码试图实现一个自定义的 `io.Writer`，但有 3 个问题：

```go
// 问题 1：方法签名错误
type BufferWriter struct {
	data []byte
}

func (b BufferWriter) Write(data []byte) {  // 缺少返回值
	b.data = append(b.data, data...)
}

// 问题 2：接收者类型导致无法修改
func (b BufferWriter) Reset() {  // 值接收者
	b.data = nil
}

// 问题 3：使用空接口但没做类型检查
func writeTo(w any, msg string) {
	w.Write([]byte(msg))  // 编译错误
}
```

<details>
<summary>点击查看答案</summary>

```go
type BufferWriter struct {
	data []byte
}

// 修复 1：实现 io.Writer 需要 (int, error) 返回值
func (b *BufferWriter) Write(data []byte) (int, error) {
	b.data = append(b.data, data...)
	return len(data), nil
}

// 修复 2：用指针接收者才能修改
func (b *BufferWriter) Reset() {
	b.data = nil
}

// 修复 3：用接口类型或做类型断言
func writeTo(w io.Writer, msg string) error {
	_, err := w.Write([]byte(msg))
	return err
}
```

</details>

### 练习 3：实现 Stringer 接口

为 `Person` 结构体实现 `fmt.Stringer` 接口，让 `fmt.Println(person)` 输出格式化信息：

```go
type Person struct {
	Name string
	Age  int
	City string
}

// 你的代码：实现 String() 方法

func main() {
	p := Person{Name: "Alice", Age: 30, City: "Taipei"}
	fmt.Println(p)  // 期望输出：Person(Alice, 30, Taipei)
}
```

<details>
<summary>点击查看答案</summary>

```go
type Person struct {
	Name string
	Age  int
	City string
}

func (p Person) String() string {
	return fmt.Sprintf("Person(%s, %d, %s)", p.Name, p.Age, p.City)
}
```

**原理**：`fmt.Stringer` 接口只要求 `String() string` 方法，实现后 `fmt` 包会自动调用。

</details>

## 故障排查 (FAQ)

### Q1: 如何判断一个类型是否实现了某个接口？

**A**: 最简单的方法是尝试赋值或传递：

```go
type MyType struct{}
func (m MyType) Read(p []byte) (int, error) { return 0, nil }

// 编译通过就表示实现了
var r io.Reader = MyType{}  // ✅ 如果报错就是没实现
```

也可以用 `var _ io.Reader = MyType{}` 做编译期检查（下划线表示不需要值）。

### Q2: 接口和泛型（Generics）应该用哪个？

**A**: 根据场景选择：

- **用接口**：当你需要运行时多态，或不同实现有不同行为时
- **用泛型**：当你需要编译时类型检查，或操作容器/算法时

```go
// 接口：运行时决定
func Process(r io.Reader) { /* ... */ }

// 泛型：编译时决定
func Map[T any](slice []T, fn func(T) T) []T { /* ... */ }
```

### Q3: 如何调试接口值的具体类型？

**A**: 用 `%T` 格式化动词：

```go
var r io.Reader = strings.NewReader("hello")
fmt.Printf("%T\n", r)  // 输出：*strings.Reader

var a any = 42
fmt.Printf("%T\n", a)  // 输出：int
```

或在调试器中查看接口的动态类型字段。

## 知识扩展 (选学)

### 接口值内部结构

接口在运行时包含两个指针：**动态类型**和**动态值**：

```
interface {
	dtype *_type   // 动态类型
	data  unsafe.Pointer  // 动态值
}
```

这解释了为什么 `var i Foo = (*T)(nil)` 不是 nil：类型是 `*T`，值是 nil。

### 接口性能开销

接口调用有轻微的运行时开销（动态分发），但在绝大多数场景下可以忽略：

```go
// 直接调用（无开销）
s := MyStruct{}
s.Method()

// 接口调用（有微小开销）
var i MyInterface = s
i.Method()
```

**基准测试**：接口调用比直接调用慢约 5-10%，但换来的是灵活性。除非在性能关键路径（如 tight loop），否则不必担心。

### 断言到具体类型

```go
func process(r io.Reader) {
	// 如果知道具体类型，可以断言获取额外方法
	if buffer, ok := r.(*bytes.Buffer); ok {
		fmt.Println("Buffer length:", buffer.Len())
	}
	
	// 或使用类型分支
	switch v := r.(type) {
	case *os.File:
		fmt.Println("File:", v.Name())
	case *bytes.Buffer:
		fmt.Println("Buffer:", v.Len())
	}
}
```

## 工业界应用

### 场景 1：可测试的 HTTP 处理器

```go
// 定义接口而非依赖具体实现
type UserRepository interface {
	GetByID(id string) (*User, error)
	Create(user *User) error
}

// 处理器依赖接口
type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// 测试时用 mock 实现
type MockUserRepo struct{}
func (m MockUserRepo) GetByID(id string) (*User, error) {
	return &User{ID: id, Name: "test"}, nil
}

func TestUserHandler(t *testing.T) {
	handler := NewUserHandler(MockUserRepo{})
	// ... 测试
}
```

### 场景 2：日志输出抽象

```go
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

// 生产环境：输出到文件
type FileLogger struct{ /* ... */ }
func (f FileLogger) Info(msg string) { /* 写入文件 */ }

// 测试环境：输出到内存
type MemoryLogger struct {
	messages []string
}
func (m *MemoryLogger) Info(msg string) {
	m.messages = append(m.messages, msg)
}

// 业务代码不关心实现
func StartServer(logger Logger) {
	logger.Info("server started")
}
```

### 场景 3：插件系统

```go
// 插件接口
type Plugin interface {
	Name() string
	Init(cfg map[string]any) error
	Execute(ctx context.Context) error
}

// 插件注册
var plugins = make(map[string]Plugin)

func Register(p Plugin) {
	plugins[p.Name()] = p
}

// 运行时加载
func RunPlugin(name string) error {
	p, ok := plugins[name]
	if !ok {
		return fmt.Errorf("plugin not found: %s", name)
	}
	return p.Execute(context.Background())
}
```

## 小结

**核心要点**：
- 隐式实现（implicit implementation）：有方法就自动实现接口
- 小接口原则：1-2 个方法的接口最灵活
- io.Writer/io.Reader 是标准库的典范
- 空接口（any）可以接收任何类型，但需要类型断言
- 接口用于解耦依赖、提高可测试性

**关键术语**：
- Implicit Implementation：隐式实现，无需 implements 关键字
- Empty Interface：空接口，`interface{}` 或`any`
- Type Assertion：类型断言，从接口恢复具体类型
- Type Switch：类型分支，根据类型执行不同逻辑
- Method Set：方法集合，类型拥有的所有方法

**下一步**：
- 学习错误处理（Error Handling），`error` 本身就是接口
- 阅读标准库 `io`、`fmt`、`sort` 包的接口设计
- 实践用接口解耦业务逻辑，编写可测试代码

## 术语表

| 英文 | 中文 | 说明 |
|------|------|------|
| Interface | 接口 | 定义行为的抽象类型 |
| Implicit Implementation | 隐式实现 | 自动实现接口，无需声明 |
| Method Set | 方法集 | 类型实现的所有方法集合 |
| Empty Interface | 空接口 | `interface{}` 或`any`，可持有任何类型 |
| Type Assertion | 类型断言 | 从接口提取具体类型 `v.(T)` |
| Type Switch | 类型开关 | 根据类型分支 `switch v := x.(type)` |
| Small Interface | 小接口 | 只有 1-2 个方法的接口 |
| Polymorphism | 多态 | 同一接口调用不同实现 |
| Decoupling | 解耦 | 减少模块间依赖 |
| Mock | 模拟 | 测试时替换真实实现 |
| Dependency Injection | 依赖注入 | 通过构造函数传入依赖 |

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/interfaces/interfaces.go)
