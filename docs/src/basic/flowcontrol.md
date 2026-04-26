# 流程控制（Flow Control）

## 开篇故事

想象你是一位餐厅经理，每天都要做各种决策："如果客人超过 8 位，就安排大包厢；如果少于 4 位，就安排吧台；否则安排普通餐桌"。你还需要重复做同样的事情："给每桌上的前菜，给每桌上的主菜"。有时候你接到电话要立刻处理（比如供应商送货），但手头上的事情要先暂停，等处理完再继续。

程序执行也是类似的。流程控制（Flow Control）就是程序的"决策系统"——它决定代码在什么时候执行、执行多少次、以及如何响应不同的情况。Go 语言的流程控制设计得非常简洁：**if/else 负责条件判断**, **for 负责所有循环**, **switch 负责多路分支**, **defer 负责延迟清理**。虽然没有花哨的语法，但它们构成了所有业务逻辑的骨架。

## 本章适合谁

- 已经会写简单 Go 程序，但想系统理解控制结构的初学者
- 从其他语言转来 Go，想知道 `for` 为什么能替代 `while` 的开发者
- 对 `defer` 执行顺序模糊，想彻底搞懂的工程师
- 想写出更清晰、更易读的分支和循环逻辑的程序员

## 你会学到什么

完成本章后，你将能够：

1. **正确使用 if/else 链**：写出互斥、有序的条件判断逻辑，避免冗余检查
2. **掌握 for 的全部用法**：用经典循环、range 遍历、条件循环处理各种场景
3. **理解 switch 的安全设计**：知道为什么 Go 默认不 fallthrough，何时使用多选一匹配
4. **解释 defer 的 LIFO 顺序**：准确预测多个 defer 的执行时机，应用于资源清理
5. **选择合适的控制结构**：根据场景选择最高效、最可读的流程控制方式

## 前置要求

- 已经安装 Go 1.24+ 开发环境
- 会写基本的 `func main()` 程序
- 理解变量、常量、基础数据类型（int, string, bool）
- 了解切片（slice）的基本概念

## 第一个例子

让我们从一个完整的温度分类器开始：

```go
package main

import "fmt"

func classifyTemperature(temp int) string {
	if temp >= 35 {
		return "炎热"
	}
	
	if temp >= 20 {
		return "温暖"
	}
	
	if temp >= 10 {
		return "凉爽"
	}
	
	return "寒冷"
}

func main() {
	fmt.Println(classifyTemperature(25)) // 输出：温暖
}
```

这个例子展示了最简单的 `if` 链：条件从上到下依次检查，一旦满足某个条件就返回，后续不再执行。这就是 Go 流程控制的哲学——**清晰的线性逻辑**，不需要花哨的嵌套。

## 原理解析

### 1. if/else：条件判断的基石

Go 的 `if` 语句有几个关键特点：

**为什么没有 `else if` 关键字？** 因为 Go 直接使用 `if` 接在 `else` 后面，形成"阶梯式"判断：

```go
if score >= 90 {
    return "优秀"
} else if score >= 60 {  // 注意：else if 是两个关键字
    return "及格"
} else {
    return "继续练习"
}
```

但在实际代码中，更推荐**早期返回（early return）**风格，减少嵌套：

```go
func classifyScore(score int) string {
	if score >= 90 {
		return "优秀"
	}

	if score >= 60 {
		return "及格"
	}

	return "继续练习"
}
```

这种写法的优势是：**每个条件都是独立的**，读者不需要在大脑里维护嵌套层级。

### 2. for：统一所有循环模式

Go 只有一个 `for` 关键字，但能表达三种循环模式：

**经典计数循环**（类似 C 的 `for`）：
```go
classicTotal := 0
for i := 1; i <= 4; i++ {
	classicTotal += i
}
// classicTotal = 10
```

**Range 遍历**（类似 Python 的 `for in`）：
```go
words := []string{"go", "is", "fun"}
rangeChars := 0
for _, word := range words {
	rangeChars += len(word)
}
// rangeChars = 7 (2 + 2 + 3)
```

注意 `for _, word` 中的下划线：Go 的 range 返回索引和值，如果不需要索引用 `_` 显式丢弃。

**条件循环**（替代 `while`）：
```go
for counter < 10 {  // 没有括号，直接写条件
	counter++
}
```

**为什么 Go 没有 while？** 因为 `for condition` 已经能完全表达 `while` 的语义，多一个关键字只会增加记忆负担。

### 3. switch：安全的多分支匹配

Go 的 `switch` 有个重要的安全特性：**默认不 fallthrough**（不自动穿透到下一个 case）：

```go
func labelDay(day string) string {
	switch day {
	case "Saturday", "Sunday":  // 多值匹配
		return "周末"
	case "Monday":
		return "新的开始"
	default:
		return "工作日"
	}
}
```

如果要显式穿透，必须写 `fallthrough` 关键字：
```go
switch grade {
case "A":
	fmt.Println("优秀")
	fallthrough  // 必须显式声明
case "B":
	fmt.Println("良好")  // A 级也会执行这行
}
```

这种设计避免了 C/Java 中常见的"忘记 break"错误。

### 4. defer：延迟执行的清理动作

`defer` 是 Go 最独特的设计之一。注册的 defer 函数会在**当前函数返回前**执行，顺序是**后进先出（LIFO）**：

```go
func collectDeferStack() (steps []string) {
	steps = append(steps, "enter")
	defer func() {
		steps = append(steps, "defer:first")
	}()
	defer func() {
		steps = append(steps, "defer:second")
	}()
	steps = append(steps, "leave")
	
	return steps
}
```

执行结果是：`enter -> leave -> defer:second -> defer:first`

**为什么会这样？** 可以想象成往箱子里放东西：先放进去的（first defer）在最下面，后放进去的（second defer）在上面。返回时从上面开始拿，所以后注册的先执行。

**实际应用场景**：
```go
func readFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()  // 无论后面是否出错，都会关闭文件
	
	data, err := io.ReadAll(f)
	// ... 处理数据
	return nil
}
```

### 5. 控制结构的组合艺术

真实代码中，这些结构经常组合使用：

```go
for _, user := range users {
	if user.IsActive {
		switch user.Role {
		case "admin":
			// ...
		case "guest":
			// ...
		}
		defer logUserAccess(user.ID)
	}
}
```

关键原则：**保持逻辑清晰**，避免过深嵌套（建议不超过 3 层）。

## 常见错误

### 错误 1：在 defer 中使用错误的变量捕获

```go
// ❌ 错误：defer 中的 i 是循环结束后的值
for i := 0; i < 3; i++ {
	defer fmt.Println(i)
}
// 输出：2, 2, 2 （不是 0, 1, 2）

// ✅ 正确：用参数捕获当前值
for i := 0; i < 3; i++ {
	defer func(val int) {
		fmt.Println(val)
	}(i)
}
// 输出：2, 1, 0 （LIFO 顺序）
```

### 错误 2：range 遍历修改原切片元素

```go
numbers := []int{1, 2, 3}

// ❌ 错误：val 是副本，不会修改原切片
for _, val := range numbers {
	val = val * 2
}
// numbers 仍然是 [1, 2, 3]

// ✅ 正确：用索引访问
for i := range numbers {
	numbers[i] = numbers[i] * 2
}
// numbers 变成 [2, 4, 6]
```

### 错误 3：switch 忘记 default 分支导致逻辑漏洞

```go
func getPriority(level string) int {
	switch level {
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	}
	return 0  // 隐式返回，容易被忽略
}

// ✅ 更好的写法：显式 default
func getPriority(level string) int {
	switch level {
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	default:
		log.Printf("unknown level: %s", level)
		return 0
	}
}
```

## 动手练习

### 练习 1：预测输出结果

```go
func mystery() string {
	steps := []string{}
	
	steps = append(steps, "A")
	defer func() {
		steps = append(steps, "D")
	}()
	
	steps = append(steps, "B")
	defer func() {
		steps = append(steps, "C")
	}()
	
	steps = append(steps, "E")
	
	return strings.Join(steps, "-")
}

fmt.Println(mystery())
// 问：输出是什么？
```

<details>
<summary>点击查看答案</summary>

```
输出：A-B-E-C-D
```

**解析**：defer 在 return 前执行（LIFO），所以 return 后、真正返回前，先执行第二个 defer（加 C），再执行第一个 defer（加 D）。

</details>

### 练习 2：修复错误代码

下面的函数有 3 个问题，请修复：

```go
func processUsers(users []string) {
	for i := 0; i <= len(users); i++ {  // 问题 1
		if users[i] == "admin" {
			fmt.Println("found admin")
		}  // 问题 2：缺少对空切片的检查
		
		switch users[i] {
		case "admin":
			fmt.Println("is admin")
		case "user":
			fmt.Println("is user")
		}  // 问题 3：没有 default
	}
}
```

<details>
<summary>点击查看答案</summary>

```go
func processUsers(users []string) {
	if len(users) == 0 {
		return
	}
	
	for i := 0; i < len(users); i++ {  // 修复 1：< 而不是 <=
		if users[i] == "admin" {
			fmt.Println("found admin")
		}
		
		switch users[i] {
		case "admin":
			fmt.Println("is admin")
		case "user":
			fmt.Println("is user")
		default:  // 修复 3：添加 default
			fmt.Println("unknown role")
		}
	}
}
```

**修复要点**：
1. `i <= len(users)` 会越界（索引最大是 `len-1`）
2. 先检查空切片避免 panic
3. 添加 default 处理未知情况

</details>

### 练习 3：编写斐波那契函数

用 `for` 循环实现斐波那契数列，返回前 n 项：

```go
func fibonacci(n int) []int {
	// 你的代码
}

// 期望输出：fibonacci(6) => [0, 1, 1, 2, 3, 5]
```

<details>
<summary>点击查看答案</summary>

```go
func fibonacci(n int) []int {
	if n <= 0 {
		return []int{}
	}
	
	if n == 1 {
		return []int{0}
	}
	
	result := []int{0, 1}
	for i := 2; i < n; i++ {
		next := result[i-1] + result[i-2]
		result = append(result, next)
	}
	
	return result
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么我的 defer 没有执行？

**A**: defer 只在函数返回时执行。如果你在 defer 之前调用了 `os.Exit()`，程序会立即退出，defer 不会执行：

```go
func badExample() {
	defer fmt.Println("不会打印")
	os.Exit(1)  // 程序立即退出
}

// ✅ 正确：在最后一行返回错误，让上层决定 Exit
func goodExample() error {
	defer fmt.Println("会打印")
	return errors.New("something failed")
}
```

### Q2: for range 中的变量为什么是副本？

**A**: range 返回的变量是迭代值的**副本**，这是 Go 的设计选择。如果需要修改原元素，必须用索引访问：

```go
// ❌ 不能修改
for _, item := range slice {
	item.Modified = true
}

// ✅ 可以修改
for i := range slice {
	slice[i].Modified = true
}
```

### Q3: switch 中如何匹配类型而不是值？

**A**: 使用类型 switch（type switch）：

```go
func inspect(value any) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("整数：%d", v)
	case string:
		return fmt.Sprintf("字符串：%s", v)
	default:
		return fmt.Sprintf("未知类型：%T", value)
	}
}
```

## 知识扩展 (选学)

### 带初始化语句的 if/switch

Go 允许在 if/switch 前加初始化语句，作用域限制在条件块内：

```go
if err := checkInput(); err != nil {
	return err
}
// err 在这里已经超出作用域

switch status := getStatus(); status {
case "active":
	// ...
}
```

这种写法让错误检查和条件判断更紧凑。

### goto（不推荐使用）

Go 支持 `goto`，但官方建议避免使用。唯一推荐场景是从深层嵌套中跳出：

```go
for ... {
	for ... {
		for ... {
			if error {
				goto cleanup
			}
		}
	}
}
cleanup:
// 清理代码
```

但在大多数情况下，**抽取成函数**是更好的选择。

### 标签（label）和循环控制

可以用标签控制外层循环的 break/continue：

```go
outerLoop:
for i := 0; i < 10; i++ {
	for j := 0; j < 10; j++ {
		if shouldStop(i, j) {
			break outerLoop  // 跳出外层循环
		}
	}
}
```

## 工业界应用

### 场景 1：HTTP 请求处理器的分层验证

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 第 1 层：基础验证
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 第 2 层：认证
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	
	// 第 3 层：业务验证
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	
	// 处理逻辑...
}
```

这种"阶梯式验证"避免了深层嵌套，每层检查后早期返回。

### 场景 2：批量处理中的 defer 资源管理

```go
func processBatch(files []string) error {
	for _, path := range files {
		if err := processSingleFile(path); err != nil {
			log.Printf("failed: %v", err)
			continue  // 不是致命错误，继续下一个
		}
	}
	return nil
}

func processSingleFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	
	// 使用文件...
	return nil
}
```

### 场景 3：状态机的 switch 实现

```go
func (o *Order) Transition(event Event) error {
	switch o.Status {
	case "pending":
		switch event {
		case EventPay:
			o.Status = "paid"
		case EventCancel:
			o.Status = "cancelled"
		default:
			return fmt.Errorf("invalid event %v for pending", event)
		}
	case "paid":
		// ...
	}
	return nil
}
```

## 小结

**核心要点**：
- `if/else` 用于条件分支，推荐早期返回风格
- `for` 统一所有循环：计数、range、条件三种模式
- `switch` 默认不 fallthrough，更安全
- `defer` 在函数返回前执行，LIFO 顺序
- 控制结构应追求清晰可读，而非炫技

**关键术语**：
- Early Return：早期返回，减少嵌套
- LIFO：后进先出（Last In, First Out）
- Fallthrough：穿透到下一个 case
- Range：Go 的遍历语法
- Type Switch：类型分支判断

**下一步**：
- 学习函数（Functions）和错误处理（Error Handling）
- 练习用 defer 管理文件、数据库连接等资源
- 阅读标准库代码，观察成熟的流程控制模式

## 术语表

| 英文 | 中文 | 说明 |
|------|------|------|
| Flow Control | 流程控制 | 决定程序执行顺序的机制 |
| Conditional Branch | 条件分支 | if/else 根据条件选择执行路径 |
| Loop | 循环 | 重复执行某段代码 |
| Range | 范围遍历 | Go 的切片/数组/Map 遍历语法 |
| Switch | 多路分支 | 多条件匹配语句 |
| Fallthrough | 穿透 | switch 中执行完一个 case 后继续下一个 |
| Defer | 延迟执行 | 注册在函数返回前执行的函数 |
| LIFO | 后进先出 | Last In, First Out，栈的执行顺序 |
| Early Return | 早期返回 | 在函数开头检查条件并提前返回 |
| Type Switch | 类型开关 | 根据值的类型进行分支判断 |

[源码](https://github.com/savechina/hello-go/blob/main/internal/basic/flowcontrol/flowcontrol.go)
