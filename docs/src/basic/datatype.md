# 基础数据类型（Data Types）

## 开篇故事

想象你在整理一个仓库。你需要不同的容器来装不同的东西：小盒子装螺丝、大箱子装工具、标签柜装文件。如果所有东西都塞进同一个大袋子，找的时候会非常混乱。

Go 的**数据类型**就是这些不同的容器——整数装数字、字符串装文本、布尔包装是/否、切片装列表、map 装键值对。选对容器，你的代码才会整洁高效。

---

## 本章适合谁

如果你想理解 Go 有哪些基本数据类型、什么时候用哪种，本章适合你。你需要理解变量声明（`var` 和 `:=`），不需要其他前置知识。

---

## 你会学到什么

完成本章后，你可以：

1. 使用整数（int）、浮点数（float64）、布尔值（bool）、字符串（string）表达核心数据
2. 使用切片（slice）表示可增长的列表，理解 `len` 和 `cap` 的区别
3. 使用映射（map）进行键值对的增删改查（CRUD）
4. 使用 `time.Time` 创建和操作时间
5. 区分"值语义"和"引用式行为"的差异

---

## 前置要求

- 理解变量声明（`var` 和 `:=`）
- 理解基本的函数调用

---

## 第一个例子

让我们从最常见的数据类型开始：

```go
count := 42         // 整数（int）
price := 19.95      // 浮点数（float64）
active := true      // 布尔值（bool）
label := "Go 1.24"  // 字符串（string）
```

**关键概念**：

- Go 会根据字面量自动推断类型（type inference）
- `42` → `int`，`19.95` → `float64`，`true` → `bool`，`"Go"` → `string`

---

## 原理解析

### 1. 整数（Integers）

Go 的整数类型分为有符号（`int`）和无符号（`uint`）：

| 类型   | 大小       | 范围                    | 常见用途       |
| ------ | ---------- | ----------------------- | -------------- |
| `int`  | 32 或 64 位 | 平台相关                | 计数、索引     |
| `int8` | 8 位       | -128 到 127             | 节省内存       |
| `int64`| 64 位      | 约 ±9×10¹⁸              | 大数、时间戳   |
| `uint` | 32 或 64 位 | 0 到最大值              | 位运算、长度   |

**建议**：日常开发直接用 `int`，除非你有明确的内存或范围需求。

### 2. 浮点数（Floats）

Go 有两种浮点类型：

- `float64` — 双精度，**日常开发默认选择**
- `float32` — 单精度，仅在内存敏感时使用

**重要**：浮点数不适合精确的货币计算！

```go
// ❌ 不好：浮点数精度问题
price := 19.95
total := price * 3  // 可能是 59.85000000000001

// ✅ 好：用整数表示分（cents）
priceCents := 1995
totalCents := priceCents * 3  // 精确的 5985
```

### 3. 布尔值（Booleans）

布尔值只有两个可能：`true` 或 `false`。

```go
isActive := true
hasPermission := false
```

**常见用途**：
- 开关标志（feature flags）
- 条件判断（if 语句）
- 状态检查（`ok` 模式）

### 4. 字符串（Strings）

Go 的字符串是**不可变的**（immutable）——一旦创建，就不能修改：

```go
name := "Hello"
// name[0] = 'h'  // ❌ 编译错误！字符串不可修改
name = "hello"   // ✅ 创建新字符串
```

**字符串拼接**：
```go
// 少量拼接：用 +
full := "Hello" + " " + "World"

// 大量拼接：用 strings.Builder
var b strings.Builder
for i := 0; i < 100; i++ {
    b.WriteString("item")
}
```

### 5. 切片（Slices）

切片是 Go 中最常用的数据结构——可以理解成"可增长的数组"：

```go
scores := []int{80, 85, 90}
scores = append(scores, 95)  // 追加元素
```

**`len` vs `cap`**：

这是切片最重要的概念：

```go
s := make([]int, 3, 5)  // len=3, cap=5
// len = 当前可见元素数
// cap = 底层数组还能容纳多少元素
```

**类比**：
> 想象一个有 5 个格子的书架（cap=5），目前只放了 3 本书（len=3）。你可以继续放书直到 5 本，超过 5 本时 Go 会自动换一个更大的书架。

### 6. 映射（Maps）

map 是键值对集合，类似其他语言的字典/哈希表：

```go
ages := map[string]int{"Alice": 20}
ages["Bob"] = 18       // 写入
age := ages["Alice"]   // 读取
delete(ages, "Bob")    // 删除
```

**检查键是否存在**：
```go
age, ok := ages["Charlie"]
if !ok {
    fmt.Println("Charlie 不在 map 中")
}
```

**重要**：如果键不存在，map 返回该类型的**零值**（zero value）：
- `int` → `0`
- `string` → `""`
- `bool` → `false`

### 7. 时间（time.Time）

`time.Time` 来自标准库 `time` 包，用于表示具体时刻：

```go
now := time.Now()                          // 当前时间
meeting := time.Date(2026, time.April, 5, 14, 30, 0, 0, time.UTC)
deadline := meeting.Add(48 * time.Hour)    // 加 48 小时
```

---

## 常见错误

### 错误 1: 未初始化的 map 导致 panic

```go
var m map[string]int
m["key"] = 1  // ❌ panic: assignment to entry in nil map
```

**修复方法**：

使用 `make` 或字面量初始化：
```go
m := make(map[string]int)      // ✅
m["key"] = 1

// 或者
m := map[string]int{"key": 1}  // ✅
```

---

### 错误 2: 切片共享底层数组

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]     // b = [2, 3]
b[0] = 99       // 修改 b
fmt.Println(a)  // ❌ a 变成了 [1, 99, 3, 4, 5]！
```

**为什么会这样？**

切片 `b` 和 `a` 共享同一个底层数组。修改 `b` 会影响 `a`。

**修复方法**：

如果需要独立副本，用 `copy`：
```go
a := []int{1, 2, 3, 4, 5}
b := make([]int, 2)
copy(b, a[1:3])  // ✅ 复制数据
b[0] = 99
fmt.Println(a)   // a 仍然是 [1, 2, 3, 4, 5]
```

---

### 错误 3: 用浮点数做精确计算

```go
price := 0.1 + 0.2
fmt.Println(price == 0.3)  // ❌ 输出 false！
```

**修复方法**：

用整数表示最小单位：
```go
priceCents := 10 + 20
fmt.Println(priceCents == 30)  // ✅ true
```

---

## 动手练习

### 练习 1: 预测输出

不运行代码，预测下面代码的输出：

```go
s := []int{1, 2, 3}
s = append(s, 4)
fmt.Println(len(s), cap(s))
```

<details>
<summary>点击查看答案</summary>

**输出**:
```
4 4
```

**解析**：
1. 初始切片 `[1, 2, 3]`，len=3, cap=3
2. `append` 追加 4，需要扩容，新 cap=4
3. 最终 len=4, cap=4

</details>

---

### 练习 2: 修复 panic

下面的代码会 panic，请修复：

```go
func main() {
    var users map[string]int
    users["Alice"] = 25  // ❌ panic
}
```

<details>
<summary>点击查看修复方法</summary>

**修复**：
```go
func main() {
    users := make(map[string]int)  // ✅ 初始化
    users["Alice"] = 25
}
```

</details>

---

### 练习 3: 切片截取

写出下面代码的输出：

```go
s := []int{10, 20, 30, 40, 50}
a := s[1:3]
b := s[2:]
fmt.Println("a:", a)
fmt.Println("b:", b)
fmt.Println("len(a):", len(a), "cap(a):", cap(a))
```

<details>
<summary>点击查看答案</summary>

**输出**:
```
a: [20 30]
b: [30 40 50]
len(a): 2 cap(a): 4
```

**解析**：
- `s[1:3]` 从索引 1 到 3（不含 3），得到 `[20, 30]`
- `s[2:]` 从索引 2 到最后，得到 `[30, 40, 50]`
- `a` 的 cap 是从起始位置到底层数组末尾：`[20, 30, 40, 50]` = 4

</details>

---

## 故障排查 (FAQ)

### Q: 什么时候用 slice，什么时候用 array？

**A**: 99% 的情况用 slice。

- **slice**（`[]int`）— 长度可变，日常开发默认选择
- **array**（`[3]int`）— 长度固定，类型中包含长度，很少直接使用

**唯一用 array 的场景**：你需要固定大小的值类型（如矩阵 `[3][3]float64`）。

---

### Q: map 的遍历顺序是固定的吗？

**A**: **不是**。Go 故意随机化 map 的遍历顺序，防止开发者依赖顺序。

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}
for k, v := range m {  // 每次运行顺序可能不同
    fmt.Println(k, v)
}
```

如果需要有序遍历，先提取 key 并排序：
```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

---

### Q: 字符串为什么是不可变的？

**A**: 三个原因：

1. **安全性** — 多个 goroutine 可以安全地读取同一字符串
2. **性能** — 不可变字符串可以被共享和缓存
3. **哈希稳定** — 字符串可以作为 map 的 key，哈希值不会改变

---

## 知识扩展 (选学)

### 值语义 vs 引用式行为

Go 中大多数类型是**值语义**——赋值时复制一份：

```go
a := 5
b := a  // b 是 a 的副本，修改 b 不影响 a
```

但 slice 和 map 有**引用式行为**——赋值时共享底层结构：

```go
a := []int{1, 2, 3}
b := a      // b 和 a 共享底层数组
b[0] = 99   // a[0] 也变成了 99
```

**规则**：
> 当你需要"修改不影响原值"时，用 `copy`（slice）或手动克隆（map）。

---

### 切片扩容策略

当 `append` 超过 cap 时，Go 会自动扩容：

- 旧 cap < 1024：新 cap = 旧 cap × 2
- 旧 cap ≥ 1024：新 cap = 旧 cap × 1.25

这意味着扩容是**指数增长**的，`append` 的平均时间复杂度是 O(1)。

---

## 工业界应用：用户配置存储

**场景**：存储和管理用户配置

```go
type UserConfig struct {
    Name     string
    Age      int
    Active   bool
    Tags     []string
    Settings map[string]string
}

func main() {
    config := UserConfig{
        Name:   "Alice",
        Age:    30,
        Active: true,
        Tags:   []string{"admin", "developer"},
        Settings: map[string]string{
            "theme": "dark",
            "lang":  "zh-CN",
        },
    }

    fmt.Printf("用户: %s, 年龄: %d, 活跃: %t\n",
        config.Name, config.Age, config.Active)
    fmt.Printf("标签: %v\n", config.Tags)
    fmt.Printf("主题: %s\n", config.Settings["theme"])
}
```

**为什么这样设计**：
- `string` 存储名称（不可变，安全）
- `int` 存储年龄（精确整数）
- `bool` 存储状态（是/否）
- `[]string` 存储标签（可增长列表）
- `map[string]string` 存储配置（键值对查找）

---

## 小结

**核心要点**：

1. **int 和 float64 是最常用的数值类型** - 日常开发直接用它们
2. **slice 是可增长的列表** - 理解 `len`（可见元素）和 `cap`（总容量）
3. **map 是键值对集合** - 使用前必须 `make` 初始化
4. **字符串是不可变的** - 修改会创建新字符串
5. **time.Time 处理时间** - 用 `time.Now()` 获取当前时间，`Add()` 做运算

**关键术语**：

- **Slice (切片)**: 可增长的有序集合，底层是数组
- **Map (映射)**: 键值对集合，类似字典
- **Zero Value (零值)**: 类型未赋值时的默认值
- **CRUD**: 创建（Create）、读取（Read）、更新（Update）、删除（Delete）
- **Value Semantics (值语义)**: 赋值时复制一份
- **Capacity (容量)**: 切片底层数组的总大小

**下一步**：

- 继续：[函数](functions.md)
- 回顾：[阶段复习](review-basic.md)

---

## 术语表

| English       | 中文     |
| ------------- | -------- |
| Integer       | 整数     |
| Float         | 浮点数   |
| Boolean       | 布尔值   |
| String        | 字符串   |
| Slice         | 切片     |
| Map           | 映射     |
| Capacity      | 容量     |
| Zero Value    | 零值     |
| Value Semantics | 值语义 |

---

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/datatype/datatype.go)
