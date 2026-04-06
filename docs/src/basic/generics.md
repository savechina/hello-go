# 泛型（Generics）

## 开篇故事

想象你在工厂流水线上工作。有一天，老板让你做一个"包装盒"的函数：给整数打包、给字符串打包、给浮点数打包。你写了三个函数：`packInt`、`packString`、`packFloat64`。第二天，老板说还要支持布尔值、自定义结构体……你意识到，这样写下去永远写不完。

泛型就是为了解决这个问题而生的。它允许你写一个"通用包装盒"，告诉编译器："我这个函数能处理任何类型，但具体是什么类型，调用时再决定。"在 Go 1.18 之前，我们只能用 `interface{}`，但那样会丢失类型安全。泛型让我们在保持类型安全的同时，实现真正的代码复用。

本章将通过真实的代码示例，带你理解泛型的核心：类型参数、类型约束，以及如何用泛型编写可复用的数据结构。

## 本章适合谁

- 你已经写过一些 Go 代码，见过 `func Foo[T any](x T) T` 这种语法，但不完全理解
- 你写过重复的函数（比如 `SumInts`、`SumFloats`），想用一种方式统一它们
- 你想理解 `comparable`、`~int` 这些约束到底有什么用
- 你打算写通用的数据结构（比如栈、队列、链表），不想为每种类型写一遍

如果你还没写过 Go 函数，建议先学习 [函数基础](./functions.md)；如果你只想用现成的泛型库，可以直接跳到 [标准库泛型示例](../advanced/stdlib-generics.md)。

## 你会学到什么

学完本章，你将能够：

1. 定义泛型函数，使用类型参数 `[T any]` 复用逻辑
2. 编写类型约束（constraints），限制 `T` 只能是某些类型
3. 理解 `comparable` 约束的用途和使用场景
4. 创建泛型类型（如 `stack[T]`），为泛型结构体编写方法
5. 使用泛型编写高阶函数（如 `mapSlice`），组合函数和类型参数

## 前置要求

在开始之前，你需要：

- **Go 1.18+**：泛型是在 Go 1.18 引入的，本章示例基于 Go 1.24
- **理解接口（interface）**：类型约束本质上是接口，你需要知道接口如何定义行为
- **理解切片（slice）**：示例中大量使用 `[]T`，你需要熟悉切片操作
- **理解函数参数**：泛型函数的参数分为"类型参数"和"普通参数"，概念上要区分

如果对这些概念不熟悉，建议先阅读：[接口](./interfaces.md)、[切片](./slices.md)、[函数](./functions.md)。

## 第一个例子

让我们从最简单、最经典的例子开始：一个能处理多种数值类型的求和函数。

### 没有泛型的时候

在泛型出现之前，你可能需要写多个版本：

```go
func sumInts(values []int) int {
    var total int
    for _, v := range values {
        total += v
    }
    return total
}

func sumFloat64(values []float64) float64 {
    var total float64
    for _, v := range values {
        total += v
    }
    return total
}
```

这两段代码逻辑完全一样，只是类型不同。如果还要支持 `int64`、`uint32`，代码量会成倍增长。

### 使用泛型

用泛型改写后，只需要一个函数：

```go
type number interface {
    ~int | ~int64 | ~float64
}

func sumValues[T number](values []T) T {
    var total T
    for _, value := range values {
        total += value
    }
    return total
}
```

调用时，编译器会自动推断类型：

```go
ints := []int{1, 2, 3}
total := sumValues(ints)  // T 被推断为 int

floats := []float64{1.5, 2.5}
avg := sumValues(floats)  // T 被推断为 float64
```

这个例子展示了泛型的核心价值：**逻辑不变，类型可变**。

## 原理解析

### 1. 类型参数（Type Parameters）

类型参数是泛型的核心。`[T number]` 里的 `T` 就像普通函数的参数，只不过它代表的是"类型"而不是"值"。

```go
func sumValues[T number](values []T) T {
    //              ^^^^^^^^  类型参数声明
    //                       ^ 返回值使用类型参数
}
```

- **声明位置**：类型参数写在函数名之后、普通参数之前，用方括号 `[]` 包裹
- **使用方式**：在函数签名中，`T` 可以出现在参数类型、返回值类型中
- **推断机制**：调用时，编译器根据传入的实参自动推断 `T` 是什么

### 2. 类型约束（Type Constraints）

类型约束限制了 `T` 可以是哪些类型。`number` 是一个接口，但它用作约束：

```go
type number interface {
    ~int | ~int64 | ~float64
}
```

这里有三个关键点：

- **并集约束**：`|` 表示"或"，`T` 可以是 `int`、`int64`、`float64` 中的任意一个
- **底层类型匹配**：`~int` 表示"底层类型是 int 的所有类型"，包括 `type myInt int` 这样的自定义类型
- **约束即接口**：约束本质是接口，可以定义方法集，也可以定义类型并集

### 3. comparable 约束

`comparable` 是 Go 内置的约束，表示"可以用 `==` 比较的类型"：

```go
func contains[T comparable](values []T, target T) bool {
    for _, value := range values {
        if value == target {  // 只有 comparable 类型才能用 ==
            return true
        }
    }
    return false
}
```

为什么需要这个约束？因为不是所有类型都能用 `==` 比较。比如切片、映射、函数类型的值不能直接比较。`comparable` 告诉编译器："放心，这个类型支持 `==` 操作。"

### 4. 泛型类型（Generic Types）

泛型不仅用于函数，还可以定义泛型结构体：

```go
type stack[T any] struct {
    items []T
}

func (s *stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, true
}
```

**关键点**：
- 结构体定义时声明类型参数 `[T any]`
- 方法接收者也要声明 `[T]`，如 `func (s *stack[T]) Push`
- 方法内部可以直接使用 `T`

### 5. 多类型参数

一个泛型函数可以有多个类型参数：

```go
func mapSlice[T any, R any](values []T, mapper func(T) R) []R {
    result := make([]R, 0, len(values))
    for _, value := range values {
        result = append(result, mapper(value))
    }
    return result
}
```

这里 `T` 是输入类型，`R` 是输出类型。调用时：

```go
numbers := []int{1, 2, 3}
doubled := mapSlice(numbers, func(n int) int { return n * 2 })
labels := mapSlice(numbers, func(n int) string { return fmt.Sprintf("%d", n) })
```

## 常见错误

### 错误 1：忘记声明类型约束

```go
// 错误
func sumValues[T](values []T) T {
    var total T
    for _, v := range values {
        total += v  // 编译错误：T 可能不支持 + 操作
    }
    return total
}
```

**修复**：添加约束，确保 `T` 是数值类型。

```go
// 正确
func sumValues[T number](values []T) T {
    var total T
    for _, v := range values {
        total += v
    }
    return total
}
```

### 错误 2：对不可比较类型使用 `comparable`

```go
// 错误
func findSlice[T comparable](slices [][]T, target []T) bool {
    for _, s := range slices {
        if s == target {  // 编译错误：切片不是 comparable 类型
            return true
        }
    }
    return false
}
```

**修复**：用 `bytes.Equal` 或自定义比较函数，而不是泛型约束。

```go
// 正确
func findIntSlice(slices [][]int, target []int) bool {
    for _, s := range slices {
        if slicesEqual(s, target) {
            return true
        }
    }
    return false
}

func slicesEqual(a, b []int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
```

### 错误 3：方法接收者忘记类型参数

```go
// 错误
type stack[T any] struct {
    items []T
}

func (s *stack) Push(item T) {  // 编译错误：未定义 T
    s.items = append(s.items, item)
}
```

**修复**：方法接收者也要声明类型参数。

```go
// 正确
func (s *stack[T]) Push(item T) {
    s.items = append(s.items, item)
}
```

## 动手练习

### 练习 1：泛型过滤器

写一个 `Filter` 函数，接受切片和谓词函数，返回满足条件的元素：

```go
func Filter[T any](values []T, predicate func(T) bool) []T {
    // 你的代码
}
```

要求：
- 输入 `[1, 2, 3, 4, 5]` 和 `func(n int) bool { return n%2 == 0 }`，输出 `[2, 4]`
- 输入 `[]string{"go", "rust", "python"}` 和 `func(s string) bool { return len(s) > 3 }`，输出 `["rust", "python"]`

<details>
<summary>参考答案</summary>

```go
func Filter[T any](values []T, predicate func(T) bool) []T {
    result := make([]T, 0)
    for _, v := range values {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}
```

</details>

### 练习 2：泛型键值对

定义一个泛型的 `Pair[T, U]` 结构体，有 `Key` 和 `Value` 两个字段，并实现 `String()` 方法：

```go
type Pair[T any, U any] struct {
    Key   T
    Value U
}

func (p Pair[T, U]) String() string {
    // 你的代码
}
```

<details>
<summary>参考答案</summary>

```go
type Pair[T any, U any] struct {
    Key   T
    Value U
}

func (p Pair[T, U]) String() string {
    return fmt.Sprintf("%v: %v", p.Key, p.Value)
}
```

</details>

### 练习 3：类型约束实践

定义一个约束 `ordered`，要求类型支持 `<`、`>` 比较，然后写一个 `Min` 函数返回切片中的最小值：

```go
type ordered interface {
    // 你的约束定义
}

func Min[T ordered](values []T) T {
    // 你的代码
}
```

<details>
<summary>参考答案</summary>

```go
type ordered interface {
    ~int | ~int64 | ~float64 | ~string
}

func Min[T ordered](values []T) T {
    if len(values) == 0 {
        var zero T
        return zero
    }
    min := values[0]
    for _, v := range values[1:] {
        if v < min {
            min = v
        }
    }
    return min
}
```

</details>

## 故障排查 (FAQ)

### Q1: 泛型和 `interface{}` 有什么区别？为什么要用泛型？

**A**: 主要区别在于类型安全：

```go
// 使用 interface{}（不推荐）
func sum(values []interface{}) interface{} {
    // 类型信息丢失，需要类型断言
    total := 0
    for _, v := range values {
        total += v.(int)  // 运行时可能 panic
    }
    return total
}

// 使用泛型（推荐）
func sum[T number](values []T) T {
    var total T
    for _, v := range values {
        total += v  // 编译期类型检查
    }
    return total
}
```

泛型的优势：
- **编译期检查**：类型错误在编译时就能发现
- **无需类型断言**：类型是已知的
- **性能更好**：编译器可以为具体类型生成优化代码

### Q2: 什么时候应该用泛型？什么时候不应该用？

**A**: 适用场景：
- 数据结构（栈、队列、链表、树）需要支持多种类型
- 算法函数（排序、查找、过滤）逻辑相同，只是类型不同
- 工具函数（如 `mapSlice`、`filter`）需要保持通用性

不适用场景：
- 只处理一种特定类型（直接用具体类型更清晰）
- 类型之间行为差异很大（用接口更符合意图）
- 代码可读性会因此下降（泛型不是炫技工具）

### Q3: `~int` 和 `int` 作为约束有什么区别？

**A**: `~int` 表示"底层类型是 int 的所有类型"，包括自定义类型：

```go
type number1 interface {
    int  // 只能是 int 本身
}

type number2 interface {
    ~int  // 可以是 int 或 type myInt int
}

type myInt int

// 使用 number1
func f1[T number1](v T) {}  // f1(myInt(5)) 编译错误

// 使用 number2
func f2[T number2](v T) {}  // f2(myInt(5)) 没问题
```

实践中，**优先使用 `~int`**，这样更灵活。

## 知识扩展 (选学)

### 1. 约束嵌入（Constraint Embedding）

约束可以像接口一样嵌入其他约束：

```go
type numeric interface {
    ~int | ~int64 | ~float64
}

type ordered interface {
    numeric  // 嵌入 numeric
    ~string  // 再加上 string
}

func Min[T ordered](values []T) T {
    // 可以使用 < > 和 + - * /
}
```

### 2. 泛型工厂函数

可以用泛型编写构造函数：

```go
func NewSlice[T any](initial ...T) []T {
    return append([]T{}, initial...)
}

func NewMap[K comparable, V any]() map[K]V {
    return make(map[K]V)
}
```

### 3. 泛型与接口的组合

泛型和接口可以结合使用。比如 `sort` 包的新泛型版本：

```go
func Sort[S ~[]E, E constraints.Ordered](x S) {
    slices.Sort(x)
}
```

这里 `S` 是切片类型，`E` 是元素类型，`constraints.Ordered` 是标准库提供的预定义约束。

### 4. 类型推断的边界

编译器会自动推断类型，但有时会失败：

```go
// 推断成功
sumValues([]int{1, 2, 3})

// 推断失败，需要显式指定
var empty []int
sumValues[int](empty)  // 必须写 [int]
```

当参数为空切片或没有参数时，通常需要显式指定类型。

## 工业界应用

### 场景：通用集合工具库

在大型 Go 项目中，经常会有一套通用的集合操作工具。比如某电商平台的订单处理系统：

```go
// 从订单中筛选出金额大于阈值的
largeOrders := Filter(orders, func(o Order) bool {
    return o.Amount > 1000
})

// 提取所有订单 ID
orderIDs := Map(orders, func(o Order) string {
    return o.ID
})

// 检查是否有已取消的订单
hasCancelled := Contains(orderStatuses, StatusCancelled)
```

这些操作如果使用泛型，代码会非常简洁：

```go
type Order struct {
    ID     string
    Amount float64
}

largeOrders := Filter(orders, func(o Order) bool { return o.Amount > 1000 })
orderIDs := Map(orders, func(o Order) string { return o.ID })
hasCancelled := Contains(orderStatuses, StatusCancelled)
```

### 场景：通用 Repository 模式

在 DDD（领域驱动设计）中，Repository 通常需要为不同实体实现类似的 CRUD 方法。使用泛型后：

```go
type Repository[T any] interface {
    Get(ctx context.Context, id string) (*T, error)
    List(ctx context.Context) ([]*T, error)
    Save(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id string) error
}

// 具体实现
type userRepo struct {
    db *sql.DB
}

func (r *userRepo) Get(ctx context.Context, id string) (*User, error) {
    // 实现
}

// 或者用泛型实现通用层
type genericRepo[T any] struct {
    db *sql.DB
}

func (r *genericRepo[T]) Get(ctx context.Context, id string) (*T, error) {
    // 通用实现
}
```

这样，对于每个实体（User、Order、Product），不需要重复编写相同的查询逻辑。

### 真实案例：标准库 `slices` 包

Go 1.21 在标准库中引入了 `slices` 包，大量使用泛型：

```go
import "cmp"
import "slices"

slices.Sort(numbers)              // 排序切片
slices.Contains(tags, "important") // 检查是否包含
idx := slices.Index(items, target) // 查找索引
slices.Reverse(data)              // 反转切片
```

这些函数都是泛型的，支持任何可比较或有序的类型。

## 小结

本章我们学习了：

1. **类型参数**：`[T any]` 让我们编写通用函数和类型
2. **类型约束**：接口形式的约束（如 `number`、`comparable`）限制类型范围
3. **泛型类型**：结构体可以是泛型的，如 `stack[T]`
4. **多类型参数**：一个函数可以有多个类型参数，如 `mapSlice[T, R]`
5. **实际应用**：集合操作、Repository 模式、标准库泛型包

关键术语：
- **类型参数（Type Parameter）**：函数或类型的泛型参数
- **类型约束（Type Constraint）**：限制类型参数范围的接口
- **Comparable**：支持 `==` 比较的类型
- **泛型类型（Generic Type）**：带有类型参数的结构体或接口

下一步建议：
- 阅读 Go 官方泛型教程：https://go.dev/tour/generics/1
- 学习标准库 `slices`、`maps` 包的源码实现
- 尝试用泛型重构你项目中的重复代码

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 泛型 | Generics | 允许类型作为参数的编程范式 |
| 类型参数 | Type Parameter | 泛型函数或类型中的类型占位符，如 `[T]` |
| 类型约束 | Type Constraint | 限制类型参数范围的接口定义 |
| Comparable | Comparable | Go 内置约束，表示可用 `==` 比较的类型 |
| 底层类型 | Underlying Type | `~T` 表示匹配所有底层类型为 T 的类型 |
| 泛型类型 | Generic Type | 带有类型参数的结构体或接口 |
| 类型推断 | Type Inference | 编译器自动确定类型参数的具体类型 |

## 相关资源

- [Go 官方泛型提案](https://go.googlesource.com/proposal/+/refs/heads/master/design/go2draft-type-parameters.md)
- [Go 1.18 泛型发布说明](https://go.dev/blog/go1.18)
- [`slices` 包文档](https://pkg.go.dev/slices)
- [`maps` 包文档](https://pkg.go.dev/maps)

[源码](../../internal/basic/generics/generics.go)
