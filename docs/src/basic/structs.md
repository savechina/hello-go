# 结构体（Structs）

## 开篇故事

想象你在组装乐高积木。单个积木块就像基础数据类型（int、string 等），它们能表达简单的信息，但无法描述复杂的事物。如果你想搭一座房子，需要把多个积木组合成墙壁、屋顶、门窗——这些组合体就是结构体。

在编程世界里，结构体（Structs）是组织数据的基本单元。一个用户（User）不只是一个字符串或数字，而是姓名、年龄、地址等多个字段的组合。一个订单（Order）包含商品列表、价格、配送信息等。结构体把这些分散的数据打包成一个有意义的整体，让代码从"零散的变量"升级为"有语义的数据模型"。

Go 的结构体设计哲学很明确：**简单组合胜于复杂继承**。没有繁琐的类层次结构，没有复杂的访问修饰符，只有清晰的字段定义和方法绑定。这种设计让代码更易读、更易测试、更易维护。

## 本章适合谁

- 已经会写基本 Go 程序，想用结构体组织数据的初学者
- 从面向对象语言（Java/Python）转来 Go，想理解"组合 vs 继承"的开发者
- 对值接收者和指针接收者区别模糊的工程师
- 想设计清晰数据模型的程序员

## 你会学到什么

完成本章后，你将能够：

1. **定义和初始化结构体**：使用字面量语法创建嵌套结构体，理解字段可见性规则
2. **区分值接收者和指针接收者**：准确判断何时用哪种接收者，避免常见陷阱
3. **使用嵌入（embedding）组合行为**：通过组合复用字段和方法，而非继承
4. **设计可维护的数据模型**：为真实业务场景设计合理的结构体和关联关系
5. **读写嵌套结构体字段**：熟练访问多层嵌套的数据，理解零值（zero value）行为

## 前置要求

- 已经安装 Go 1.24+ 开发环境
- 理解变量、函数、包的基本概念
- 了解指针的基础知识（什么是指针、如何取地址）
- 知道什么是面向对象编程（类、对象、方法）

## 第一个例子

让我们从一个简单的员工档案系统开始：

```go
package main

import "fmt"

type Address struct {
	City string
}

type Profile struct {
	Name    string
	Age     int
	Address Address
}

func main() {
	p := Profile{
		Name: "Alice",
		Age:  30,
		Address: Address{
			City: "Taipei",
		},
	}
	
	fmt.Printf("%s lives in %s\n", p.Name, p.Address.City)
	// 输出：Alice lives in Taipei
}
```

这个例子展示了结构体的核心要素：**定义类型**、**嵌套字段**、**字面量初始化**、**访问字段**。

## 原理解析

### 1. 结构体定义与字段可见性

Go 的结构体由关键字 `type` 和 `struct` 定义：

```go
type Person struct {
	Name string  // 导出字段（大写开头）
	age  int     // 未导出字段（小写开头）
}
```

**为什么字段名大小写这么重要？** 在 Go 中，大写字母开头的标识符是**导出（exported）**的，其他包可以访问；小写字母开头是**未导出（unexported）**的，只能在包内访问。这是 Go 的封装机制——没有 `public`/`private` 关键字，只有大小写规则。

### 2. 初始化：字面量 vs 构造函数

Go 没有构造函数（constructor），但可以用**字面量（literal）**或**工厂函数（factory function）**初始化：

```go
// 方式 1：结构体字面量
p1 := Person{Name: "Bob", age: 25}

// 方式 2：工厂函数（推荐用于复杂初始化）
func NewPerson(name string, age int) *Person {
	return &Person{
		Name: name,
		age:  age,
	}
}

p2 := NewPerson("Carol", 27)
```

**工厂函数的优势**：
- 可以在创建时做参数校验
- 可以设置默认值
- 可以返回接口类型而非具体类型

### 3. 方法：值接收者 vs 指针接收者

这是 Go 结构体最重要的概念之一：

```go
type Counter struct {
	Value int
}

// 值接收者（value receiver）：修改的是副本
func (c Counter) IncByVal() {
	c.Value++  // 不影响原对象
}

// 指针接收者（pointer receiver）：修改的是原对象
func (c *Counter) IncByPtr() {
	c.Value++  // 影响原对象
}
```

**如何选择？** 遵循以下规则：

- 如果方法**需要修改字段**，用指针接收者
- 如果结构体**很大**（避免拷贝开销），用指针接收者
- 如果方法是**只读的**，用值接收者
- **一致性原则**：同一个类型的方法应该统一用值或指针接收者

### 4. 嵌入（embedding）：Go 式的组合

Go 没有传统继承（inheritance），但可以通过**嵌入（embedding）**实现行为复用：

```go
type Address struct {
	City string
}

type Employee struct {
	Employee   // 嵌入，没有字段名
	Department string
	Title      string
}

// 可以直接访问嵌入类型的字段
e := Employee{
	Employee: Employee{
		City: "Kaohsiung",
	},
	Department: "Platform",
	Title:      "Engineer",
}

fmt.Println(e.City)  // 直接访问，不需要 e.Employee.City
```

**嵌入 vs 继承的关键区别**：

| 特性 | 嵌入（Go） | 继承（Java/Python） |
|------|----------|------------------|
| 关键字 | 无（直接写类型名） | extends/implements |
| 多继承 | 支持（可嵌入多个类型） | 单继承为主 |
| 父类感知 | 否（嵌入类型不知道被嵌入） | 是（父类知道子类） |
| 类型转换 | 不能把子类型当父类型用 | 支持向上转型 |

### 5. 方法提升（Method Promotion）

嵌入的不只是字段，还有方法：

```go
type Greeter struct{}

func (g Greeter) SayHello() string {
	return "hello"
}

type Robot struct {
	Greeter  // 嵌入
	Series   string
}

r := Robot{Series: "R2"}
fmt.Println(r.SayHello())  // 输出：hello（方法被"提升"到 Robot）
```

这让你可以像继承一样使用嵌入类型的方法，但本质上仍是组合。

## 常见错误

### 错误 1：混淆值接收者和指针接收者

```go
type Person struct {
	Name string
	Age  int
}

// ❌ 错误：想要修改年龄但用了值接收者
func (p Person) HaveBirthday() {
	p.Age++  // 修改的是副本，原对象不变
}

// ✅ 正确：用指针接收者
func (p *Person) HaveBirthday() {
	p.Age++  // 修改原对象
}

func main() {
	p := Person{Name: "Bob", Age: 27}
	p.HaveBirthday()
	fmt.Println(p.Age)  // ❌ 输出 27（没变），✅ 输出 28
}
```

### 错误 2：嵌入后误解字段访问优先级

```go
type A struct {
	Name string
}

type B struct {
	Name string
}

type C struct {
	A
	B
}

c := C{}
c.Name = "test"  // ❌ 编译错误：ambiguous selector c.Name

// ✅ 正确：显式指定
c.A.Name = "test"
c.B.Name = "test"
```

当嵌入的多个类型有同名字段时，必须显式指定访问哪个。

### 错误 3：使用未初始化嵌套结构体导致 panic

```go
type Profile struct {
	Name    string
	Address Address
}

p := Profile{Name: "Alice"}
// p.Address 是零值（Address 结构体的零值，City 为空字符串）

// ❌ 不会 panic，但可能不是期望的行为
fmt.Println(p.Address.City)  // 输出空字符串

// ❌ 如果 Address 是指针类型会 panic
type Profile2 struct {
	Name    string
	Address *Address
}

p2 := Profile2{Name: "Alice"}
// fmt.Println(p2.Address.City)  // ❌ panic: nil pointer dereference

// ✅ 正确：初始化嵌套指针
p2.Address = &Address{City: "Taipei"}
```

## 动手练习

### 练习 1：预测输出结果

```go
type Profile struct {
	Name string
	Age  int
}

func (p Profile) summary() string {
	return fmt.Sprintf("%s is %d years old", p.Name, p.Age)
}

func (p *Profile) haveBirthday() {
	p.Age++
}

func main() {
	p := Profile{Name: "Alice", Age: 30}
	p.haveBirthday()
	fmt.Println(p.summary())
}
// 问：输出是什么？
```

<details>
<summary>点击查看答案</summary>

```
输出：Alice is 31 years old
```

**解析**：`haveBirthday()` 使用指针接收者，真正修改了 `p.Age`，所以 `summary()` 看到的是更新后的年龄。

</details>

### 练习 2：修复错误代码

下面的代码有 3 个问题，请修复：

```go
type Employee struct {
	Name       string
	Department string
	Title      string
}

// 问题 1：接收者类型错误
func (e Employee) Promote(newTitle string) {
	e.Title = newTitle
}

// 问题 2：嵌入语法错误
type Manager struct {
	Employee  // 缺少字段名但语法不对
	TeamSize int
}

func main() {
	m := Manager{
		Employee: Employee{
			Name:       "Carol",
			Department: "Engineering",
			Title:      "Engineer",
		},
		TeamSize: 5,
	}
	
	// 问题 3：调用方法后原对象没变
	m.Promote("Senior Engineer")
	fmt.Println(m.Title)  // 期望："Senior Engineer"，实际还是 "Engineer"
}
```

<details>
<summary>点击查看答案</summary>

```go
type Employee struct {
	Name       string
	Department string
	Title      string
}

// 修复 1：用指针接收者
func (e *Employee) Promote(newTitle string) {
	e.Title = newTitle
}

// 修复 2：嵌入语法原本是正确的，这里无需修改
type Manager struct {
	Employee  // 正确：匿名嵌入
	TeamSize int
}

func main() {
	m := Manager{
		Employee: Employee{
			Name:       "Carol",
			Department: "Engineering",
			Title:      "Engineer",
		},
		TeamSize: 5,
	}
	
	m.Promote("Senior Engineer")
	fmt.Println(m.Title)  // ✅ 输出：Senior Engineer
}
```

**要点**：当方法需要修改字段时，必须使用指针接收者。

</details>

### 练习 3：设计图书管理系统

定义三个结构体：`Author`（作者）、`Book`（图书）、`Library`（图书馆），满足：
- `Author` 包含姓名和出生年份
- `Book` 嵌入了 `Author`，包含书名和出版年份
- `Library` 包含图书馆名称和图书记录（切片）
- 为 `Book` 编写方法 `GetAuthorAgeAtPublication()` 计算作者出版书时的年龄

```go
// 你的代码
```

<details>
<summary>点击查看答案</summary>

```go
type Author struct {
	Name        string
	BirthYear   int
}

type Book struct {
	Author       // 嵌入
	Title        string
	PublishYear  int
}

type Library struct {
	Name   string
	Books  []Book
}

func (b Book) GetAuthorAgeAtPublication() int {
	return b.PublishYear - b.BirthYear
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么我的方法修改不了字段？

**A**: 99% 的情况是用了值接收者而不是指针接收者。检查方法签名：

```go
// ❌ 不能修改
func (s State) Update() { /* ... */ }

// ✅ 可以修改
func (s *State) Update() { /* ... */ }
```

**经验法则**：如果方法名暗示"改变"（Add、Update、Delete、Inc、Dec），几乎总是需要指针接收者。

### Q2: 如何判断结构体字段是否导出？

**A**: 看首字母大小写：

```go
type Config struct {
	Host string  // 导出：其他包可访问
	port int     // 未导出：只能在包内访问
}
```

**提示**：如果需要在包外访问但又不想暴露字段，提供 Getter/Setter 方法：

```go
func (c *Config) GetPort() int { return c.port }
```

### Q3: 嵌入和"有一个字段"有什么区别？

**A**: 直接看代码：

```go
// 嵌入（embedding）
type A struct {
	Name string
}

type B struct {
	A  // 嵌入
}

b := B{}
b.Name = "test"  // ✅ 可以直接访问

// "有一个"（has-a）
type C struct {
	a A  // 有名字段
}

c := C{}
c.a.Name = "test"  // 必须通过字段名访问
```

**选择建议**：如果是"is-a"关系（Manager is an Employee），用嵌入；如果是"has-a"关系（Car has an Engine），用有名字段。

## 知识扩展 (选学)

### 结构体标签（Struct Tags）

用于给字段添加元数据，常用于 JSON 序列化、数据库映射：

```go
type User struct {
	ID    int    `json:"id" db:"user_id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email,omitempty"`
}
```

运行时可以用 `reflect` 包读取标签：

```go
t := reflect.TypeOf(user)
field, _ := t.FieldByName("Email")
fmt.Println(field.Tag.Get("json"))  // 输出：email,omitempty
```

### 零值（Zero Value）初始化

结构体字段有默认的零值，可以不显式初始化：

```go
type Config struct {
	Host string  // 零值：""
	Port int     // 零值：0
	SSL  bool    // 零值：false
}

c := Config{}
fmt.Println(c.Host)  // 空字符串
fmt.Println(c.Port)  // 0
```

**最佳实践**：依赖零值可以简化代码：

```go
// ✅ 利用零值：false 表示不需要 SSL
type DialConfig struct {
	UseSSL bool
}

cfg := DialConfig{}  // 默认不使用 SSL
```

### 结构体比较

Go 允许直接用 `==` 比较结构体，前提是所有字段都是可比较的：

```go
type Point struct {
	X int
	Y int
}

p1 := Point{1, 2}
p2 := Point{1, 2}
fmt.Println(p1 == p2)  // 输出：true
```

但如果有切片、Map、函数类型的字段，则不能直接比较。

## 工业界应用

### 场景 1：API 请求/响应模型

```go
// 请求体
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"min=0"`
}

// 响应体
type CreateUserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// 使用工厂函数确保必填字段
func NewCreateUserRequest(name, email string, age int) *CreateUserRequest {
	return &CreateUserRequest{
		Name:  name,
		Email: email,
		Age:   age,
	}
}
```

### 场景 2：配置管理中的嵌套结构

```go
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Database DatabaseConfig `yaml:"database"`
	LogLevel string         `yaml:"log_level"`
}

// 从 YAML 文件加载
func LoadConfig(path string) (*ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	var cfg ServerConfig
	return &cfg, yaml.Unmarshal(data, &cfg)
}
```

### 场景 3：领域模型组合

```go
// 基础信息嵌入
type Auditable struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 订单模型
type Order struct {
	Auditable  // 嵌入审计字段
	ID         string
	CustomerID string
	Items      []OrderItem
	Status     OrderStatus
}

func (o *Order) AddItem(item OrderItem) {
	o.Items = append(o.Items, item)
	o.UpdatedAt = time.Now()
}
```

## 小结

**核心要点**：
- 结构体是 Go 组织数据的基本单元，用 `type X struct { ... }` 定义
- 字段首字母大小写决定可见性（导出/未导出）
- 指针接收者修改原对象，值接收者修改副本
- 嵌入（embedding）是 Go 式的组合，不是继承
- 零值初始化可简化代码，但要理解默认值含义

**关键术语**：
- Struct Literal：结构体字面量，初始化语法
- Receiver：方法接收者（值或指针）
- Embedding：嵌入，Go 的组合机制
- Exported Field：导出字段（大写开头）
- Zero Value：零值，类型的默认值
- Method Promotion：方法提升，嵌入类型的方法可直接调用

**下一步**：
- 学习接口（Interfaces），理解 Go 的隐式实现
- 练习设计合理的领域模型
- 阅读标准库源码，观察结构体最佳实践

## 术语表

| 英文 | 中文 | 说明 |
|------|------|------|
| Struct | 结构体 | 组合多个字段的数据类型 |
| Field | 字段 | 结构体的成员变量 |
| Method | 方法 | 绑定到结构体的函数 |
| Receiver | 接收者 | 方法绑定的目标（值或指针） |
| Value Receiver | 值接收者 | 方法接收结构体副本 |
| Pointer Receiver | 指针接收者 | 方法接收结构体指针，可修改原对象 |
| Embedding | 嵌入 | Go 的组合机制，类似继承但本质不同 |
| Struct Literal | 结构体字面量 | 初始化结构体的语法 |
| Exported | 导出的 | 大写字母开头，其他包可访问 |
| Unexported | 未导出的 | 小写字母开头，包内私有 |
| Zero Value | 零值 | 类型的默认值（int=0, string=""等） |
| Factory Function | 工厂函数 | 返回结构体指针的构造函数 |

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/structs/structs.go)
