# 反射（Reflection）

## 开篇故事

想象你是一位图书管理员，每天的工作是整理成千上万本书。有一天，老板要求你开发一个系统：不管送来的是什么书，系统都能自动识别书的类型（小说、科普、历史）、提取作者信息、检查书籍标签，甚至能按照特定指令把书放到正确的位置。

但你面临一个挑战：书的种类太多，你无法为每一种书都写一个专门的处理函数。这时候，你需要一种"通用阅读能力"——能够在拿到任何一本书时，实时地检查它的属性，然后根据这些信息决定如何处理。

Go 语言的反射（reflection）机制就是这种"通用阅读能力"。它允许程序在运行时（runtime）动态地检查任何变量的类型（type）、值（value）、字段（field）、方法（method）以及结构体标签（struct tag）。就像图书管理员学会了快速扫描任何书籍并提取关键信息的能力。

但反射也是一把双刃剑。使用得当，它能帮你构建灵活的框架、序列化工具、ORM 系统；滥用反射，则会让代码变得难以理解、调试困难、性能下降。本章的目标是帮你掌握反射的正确打开方式。

## 本章适合谁

- ✅ 已经掌握 Go 基础语法（结构体、接口、方法）的学习者
- ✅ 对框架原理感兴趣，想理解"为什么结构体声明就能自动序列化"的开发者
- ✅ 需要编写通用工具函数（如配置加载、数据校验、日志格式化）的工程师
- ✅ 准备学习 GORM、encoding/json、validator 等库源码的技术人员

如果你还没有写过结构体（struct）或接口（interface），建议先完成基础章节再回来学习本章。

## 你会学到什么

学完本章后，你将能够：

1. **区分 Type 与 Value**：理解 `reflect.Type` 和 `reflect.Value` 的核心差异，知道何时使用哪个
2. **解析结构体标签**：读取和处理 struct tag，理解 ORM 和 JSON 序列化背后的原理
3. **动态调用方法**：在运行时通过方法名调用函数，了解插件系统的工作机制
4. **安全使用反射**：掌握反射的边界和陷阱，知道什么时候不应该使用反射
5. **编写通用工具**：基于反射实现配置校验、元数据提取等实用功能

## 前置要求

在开始本章之前，请确保你已经掌握：

- Go 基础语法（变量、函数、控制流）
- 结构体（struct）定义和使用
- 接口（interface）基本概念
- 指针（pointer）的基本操作
- 错误处理（error handling）

如果对上述概念还不熟悉，建议先复习基础章节。

## 第一个例子

让我们从最简单的反射使用场景开始：查看一个变量的类型和值。

```go
package main

import (
	"fmt"
	"reflect"
)

type taggedUser struct {
	Name  string `json:"name" db:"user_name"`
	Level int    `json:"level" db:"user_level"`
}

func main() {
	user := taggedUser{Name: "gopher", Level: 3}
	
	// 获取类型信息
	typ := reflect.TypeOf(user)
	fmt.Printf("类型名称：%s\n", typ.Name())
	fmt.Printf("类型种类：%s\n", typ.Kind())
	
	// 获取值信息
	val := reflect.ValueOf(user)
	fmt.Printf("实际值：%v\n", val.Interface())
	
	// 组合描述
	fmt.Printf("完整描述：type=%s kind=%s value=%v\n", 
		typ.String(), val.Kind(), val.Interface())
}
```

运行结果：
```
类型名称：taggedUser
类型种类：struct
实际值：{gopher 3}
完整描述：type=main.taggedUser kind=struct value={gopher 3}
```

这个例子展示了反射最基础的用途：在运行时获取类型和值的描述。你可能会问："我直接打印不就行了吗？为什么要用反射？"

关键在于**通用性**。当你编写一个需要处理任意类型输入的函数时（比如日志记录器、序列化器），你无法预知输入是什么类型，这时反射就成为必要工具。

## 原理解析

### 概念 1：reflect.Type 与 reflect.Value

Go 反射的两大基石是 `reflect.Type` 和 `reflect.Value`：

| 特性 | reflect.Type | reflect.Value |
|------|--------------|---------------|
| 关注点 | "这是什么类型" | "这个值是什么" |
| 获取方式 | `reflect.TypeOf(x)` | `reflect.ValueOf(x)` |
| 典型用途 | 获取类型名、字段、方法、标签 | 获取/设置值、调用方法 |
| 零值检查 | `typ == nil` | `!val.IsValid()` |

理解这个区别很重要：Type 描述的是"模具"，Value 描述的是"用模具做出来的东西"。

### 概念 2：Kind 与 Type 的区别

```go
type MyInt int

var a int
var b MyInt

fmt.Println(reflect.TypeOf(a))  // int
fmt.Println(reflect.TypeOf(b))  // main.MyInt
fmt.Println(reflect.ValueOf(a).Kind())  // int
fmt.Println(reflect.ValueOf(b).Kind())  // int
```

`Type` 能看到自定义类型名（MyInt），而 `Kind` 看到的是底层基础类型（int）。在处理 switch 判断时，通常使用 `Kind()` 做分类。

### 概念 3：结构体标签（Struct Tag）

结构体标签是 Go 反射最实用的部分之一：

```go
type taggedUser struct {
	Name  string `json:"name" db:"user_name"`
	Level int    `json:"level" db:"user_level"`
}

// 读取标签
typ := reflect.TypeOf(taggedUser{})
field, _ := typ.FieldByName("Name")
fmt.Println(field.Tag.Get("json"))  // name
fmt.Println(field.Tag.Get("db"))    // user_name
```

标签本质上是**元数据（metadata）**，不会自动生效。必须有代码主动读取标签并执行相应逻辑。这就是为什么 encoding/json 能根据 `json:"name"` 自动序列化字段，GORM 能根据 `db:"user_name"` 映射数据库列名。

### 概念 4：动态方法调用

反射允许在运行时通过方法名调用函数：

```go
type greeter struct {
	Prefix string
}

func (g greeter) Greet(name string) string {
	return fmt.Sprintf("%s, %s", g.Prefix, name)
}

// 动态调用
g := greeter{Prefix: "hello"}
method := reflect.ValueOf(g).MethodByName("Greet")
args := []reflect.Value{reflect.ValueOf("Go")}
result := method.Call(args)
fmt.Println(result[0].Interface())  // hello, Go
```

这种能力适合构建**插件系统**、**命令路由器**、**通用测试工具**。但代价是失去编译期检查——方法名写错只有在运行时才会暴露。

### 概念 5：指针处理

反射处理指针时需要特别小心：

```go
type User struct {
	Name string
}

u := &User{Name: "gopher"}

// 获取指针的类型
typ := reflect.TypeOf(u)
fmt.Println(typ.Kind())  // ptr

// 获取指针指向的类型
if typ.Kind() == reflect.Ptr {
	typ = typ.Elem()
	fmt.Println(typ.Name())  // User
}

// 获取指针指向的值
val := reflect.ValueOf(u)
if val.Kind() == reflect.Ptr {
	val = val.Elem()
	fmt.Println(val.FieldByName("Name"))  // gopher
}
```

很多反射函数都要求先判断是否是 Pointer，然后通过 `Elem()` 获取实际内容。

## 常见错误

### 错误 1：不检查 Kind 就直接处理

```go
// ❌ 错误示例
func process(input any) {
	typ := reflect.TypeOf(input)
	// 直接假设 input 是 struct
	for i := 0; i < typ.NumField(); i++ {  // 如果 input 不是 struct 会 panic！
		field := typ.Field(i)
		fmt.Println(field.Name)
	}
}

// ✅ 正确示例
func process(input any) {
	typ := reflect.TypeOf(input)
	if typ == nil {
		return
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return  // 或者返回错误
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fmt.Println(field.Name)
	}
}
```

### 错误 2：混淆 Type 和 Value 的零值检查

```go
// ❌ 错误示例
func describe(input any) {
	typ := reflect.TypeOf(input)
	val := reflect.ValueOf(input)
	if typ == nil || val == nil {  // Value 没有 nil 概念
		return
	}
}

// ✅ 正确示例
func describe(input any) {
	typ := reflect.TypeOf(input)
	val := reflect.ValueOf(input)
	if typ == nil || !val.IsValid() {  // Value 用 IsValid() 检查
		return
	}
}
```

### 错误 3：反射调用时参数数量或类型不匹配

```go
// ❌ 错误示例
func callMethod(target any, method string, args ...string) string {
	val := reflect.ValueOf(target)
	selected := val.MethodByName(method)
	// 直接调用，不检查方法是否存在
	result := selected.Call(nil)  // 如果方法需要参数会 panic！
	return fmt.Sprint(result[0].Interface())
}

// ✅ 正确示例
func callMethod(target any, method string, args ...string) string {
	val := reflect.ValueOf(target)
	selected := val.MethodByName(method)
	if !selected.IsValid() {
		return "method not found"
	}
	
	// 构建正确的参数
	inputs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inputs[i] = reflect.ValueOf(arg)
	}
	
	result := selected.Call(inputs)
	if len(result) == 0 {
		return "no result"
	}
	return fmt.Sprint(result[0].Interface())
}
```

## 动手练习

### 练习 1：实现一个简单的字段提取器

编写一个函数 `ExtractFields(input any) []string`，输入任意结构体，返回所有字段名的切片。

**提示**：使用 `reflect.TypeOf()` 获取类型，然后遍历字段。

<details>
<summary>参考答案</summary>

```go
func ExtractFields(input any) []string {
	typ := reflect.TypeOf(input)
	if typ == nil {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}

	fields := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		fields = append(fields, typ.Field(i).Name)
	}
	return fields
}
```

</details>

### 练习 2：读取所有 JSON 标签

编写一个函数 `GetJSONTags(input any) map[string]string`，返回字段名到 JSON 标签的映射。

**提示**：使用 `field.Tag.Get("json")` 读取标签。

<details>
<summary>参考答案</summary>

```go
func GetJSONTags(input any) map[string]string {
	typ := reflect.TypeOf(input)
	if typ == nil || typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ == nil || typ.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]string)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			result[field.Name] = jsonTag
		}
	}
	return result
}
```

</details>

### 练习 3：安全的类型描述函数

实现本章源码中的 `describeValue` 函数，要求处理所有边界情况（nil、指针、非结构体等）。

<details>
<summary>参考答案</summary>

```go
func describeValue(input any) string {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	
	if !val.IsValid() || typ == nil {
		return "invalid value"
	}

	return fmt.Sprintf("type=%s kind=%s value=%v", 
		typ.String(), val.Kind(), val.Interface())
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么反射代码比直接写类型代码慢？

**答**：反射需要在运行时进行类型查询、方法查找、参数验证等操作，这些都会带来额外的 CPU 开销。此外，反射调用通常会绕过编译期优化。建议：

- 性能敏感路径避免反射
- 缓存 `reflect.Type` 结果（类型不会变化）
- 能用接口解决的场景优先用接口

### Q2: 反射会破坏类型安全吗？

**答**：是的，这是反射的代价。反射调用中的方法名错误、参数类型不匹配等问题只有在运行时才会暴露。降低风险的方法：

- 为反射函数编写充分的单元测试
- 在函数入口处做严格的类型检查
- 返回清晰的错误信息而非 panic

### Q3: 什么时候应该使用反射？

**答**：反射适合以下场景：

- ✅ 编写框架代码（ORM、序列化、配置加载）
- ✅ 实现通用工具函数（日志格式化、数据校验）
- ✅ 处理未知类型的输入（插件系统）
- ❌ 普通业务逻辑（应该用接口和显式类型）
- ❌ 性能敏感的热点代码

## 知识扩展 (选学)

### 扩展 1：reflect.DeepEqual 的原理

`reflect.DeepEqual` 是 testing 包中常用的比较函数，它能递归比较两个任意类型的值是否相等。理解其原理有助于编写更好的测试代码。

### 扩展 2：自定义反射行为

Go 允许通过实现特定接口来影响反射行为，例如 `encoding.TextMarshaler` 接口会影响 json 包的序列化方式。

### 扩展 3：unsafe 包与反射

`unsafe` 包提供了更底层的内存操作能力，可以与反射配合使用实现零拷贝转换。但这是高级主题，需要深入理解 Go 内存模型。

### 扩展 4：代码生成替代反射

许多现代 Go 项目使用代码生成（如 go generate）来替代反射，在编译期生成类型安全代码，同时保持灵活性。

## 工业界应用

### 场景：配置校验系统

某公司的微服务平台需要支持多种服务配置，每个服务的配置字段不同，但都需要验证必填字段、最小值、格式等规则。

**传统方案**：为每个配置类型手写校验逻辑，代码重复且容易遗漏。

**反射方案**：

```go
type ServiceConfig struct {
	ServiceName string `json:"service_name" required:"true"`
	Port        int    `json:"port" min:"1" max:"65535"`
	Host        string `json:"host" required:"true" format:"hostname"`
}

func ValidateStruct(input any) error {
	val := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	
	// 处理指针
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}
	
	// 遍历字段检查标签
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		
		// 检查 required
		if field.Tag.Get("required") == "true" && value.IsZero() {
			return fmt.Errorf("%s is required", field.Name)
		}
		
		// 检查 min/max
		// ...
	}
	return nil
}
```

这种模式被广泛应用于配置框架、表单验证、API 参数校验等场景。

## 小结

本章介绍了 Go 反射机制的核心概念和实践技巧。让我们回顾关键要点：

### 核心概念
- `reflect.Type`：描述类型信息（名称、字段、方法、标签）
- `reflect.Value`：描述具体值（可获取/设置值、调用方法）
- `Kind`：底层类型分类（struct、int、ptr 等）
- `Struct Tag`：元数据，需主动读取才生效

### 最佳实践
1. 使用反射前始终检查 Kind
2. 正确处理指针（使用 Elem()）
3. 为反射调用提供充分的错误处理
4. 避免在性能敏感路径使用反射

### 下一步
- 阅读 encoding/json 源码理解序列化实现
- 学习 GORM 源码理解 ORM 框架设计
- 尝试编写自己的配置校验工具

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 反射 | Reflection | 运行时检查类型和值的能力 |
| 类型 | Type | 描述数据的种类和结构 |
| 值 | Value | 具体的数据内容 |
| 种类 | Kind | 类型的底层分类（struct、int、ptr 等） |
| 结构体标签 | Struct Tag | 结构体字段的元数据注释 |
| 元数据 | Metadata | 描述数据的数据 |
| 动态调用 | Dynamic Invocation | 运行时通过名称调用方法 |
| 编译期检查 | Compile-time Check | 编译时验证代码正确性 |
| 零值 | Zero Value | Go 类型的默认初始值 |
| 指针 | Pointer | 存储内存地址的变量 |

## 源码

完整示例代码位于：[internal/advance/reflection/reflection.go](https://github.com/savechina/hello-go/blob/main/internal/advance/reflection/reflection.go)
