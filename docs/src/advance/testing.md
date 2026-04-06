# 测试（Testing）

## 开篇故事

想象你是一位软件工程师，接手了一个关键项目。代码已经运行了三年，没有人敢修改。你想优化一个性能瓶颈，但每次改动都会导致其他功能莫名其妙地崩溃。为什么？

因为没有测试。

另一位同事的情况完全不同：她要重构整个核心模块，自信地写了一个下午的代码。然后运行 `go test`，十分钟后，三百个测试全部通过。她提交了代码，晚上安心睡觉。

区别在哪里？**可验证性（verifiability）**。

测试不是"写完代码后额外做的作业"，而是"让代码可信的必要条件"。没有测试的代码就像没有刹车的车——也许能开，但没人敢加速。

Go 语言天生为测试设计：内置 `testing` 包、表驱动测试、基准测试、模糊测试，无需额外框架就能写出专业的测试代码。本章带你掌握 Go 测试的三大支柱：表驱动测试、基准测试、模糊测试。

## 本章适合谁

- ✅ 已掌握 Go 基础语法（函数、结构体、切片）的开发者
- ✅ 想学习如何编写可维护测试的工程师
- ✅ 遇到"代码不敢改"困境的技术人员
- ✅ 准备构建企业级 Go 项目的团队

即使你是测试新手，只要理解基础语法就能跟上本章内容。

## 你会学到什么

学完本章后，你将能够：

1. **编写表驱动测试**：用统一结构覆盖多种测试场景
2. **设计和运行基准测试**：比较不同实现的性能差异
3. **理解模糊测试基础**：用随机输入探索边界情况
4. **组织测试代码**：遵循 Go 社区最佳实践
5. **培养测试思维**：从"能跑就行"到"可验证设计"

## 前置要求

在开始本章之前，请确保你已经掌握：

- Go 基础语法（函数、切片、map）
- 基本的错误处理
- 命令行运行 Go 程序
- 对单元测试有基本概念（可选）

## 第一个例子

让我们从最简单的测试开始：测试一个成绩评级函数。

```go
// 被测函数
func gradeLabel(score int) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 60:
		return "pass"
	default:
		return "retry"
	}
}

// 测试函数（*_test.go 文件）
func TestGradeLabel(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{score: 95, want: "excellent"},
		{score: 75, want: "pass"},
		{score: 50, want: "retry"},
		{score: 90, want: "excellent"},  // 边界值
		{score: 60, want: "pass"},       // 边界值
		{score: 59, want: "retry"},      // 边界值
	}

	for _, item := range cases {
		got := gradeLabel(item.score)
		if got != item.want {
			t.Errorf("score=%d: want %q, got %q", 
				item.score, item.want, got)
		}
	}
}
```

运行测试：
```bash
$ go test -v
=== RUN   TestGradeLabel
--- PASS: TestGradeLabel (0.00s)
PASS
```

这个例子展示了**表驱动测试（table-driven test）**的核心模式：

1. 定义测试用例切片（slice of test cases）
2. 循环执行每个用例
3. 比较实际输出和期望输出

## 原理解析

### 概念 1：表驱动测试（Table-Driven Test）

表驱动测试是 Go 社区最推崇的测试组织方式：

```go
// 传统方式：一个场景一个函数
func TestGradeExcellent(t *testing.T) { /* ... */ }
func TestGradePass(t *testing.T) { /* ... */ }
func TestGradeRetry(t *testing.T) { /* ... */ }

// 表驱动方式：数据驱动
func TestGradeLabel(t *testing.T) {
	cases := []struct {
		name  string  // 可选：用例名称
		score int
		want  string
	}{
		{name: "excellent boundary", score: 95, want: "excellent"},
		{name: "pass case", score: 75, want: "pass"},
		{name: "retry case", score: 50, want: "retry"},
		// 新增场景只需添加一行数据
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {  // 使用 t.Run 显示子测试名
			got := gradeLabel(tc.score)
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}
```

**优势**：

- ✅ **结构统一**：所有用例遵循相同模式
- ✅ **易于扩展**：新增场景只需添加数据
- ✅ **便于审查**：一眼看清覆盖了多少情况
- ✅ **减少重复**：避免复制粘贴测试逻辑

### 概念 2：基准测试（Benchmark）

基准测试用于测量函数性能：

```go
// 待比较的两种字符串拼接方式
func joinWordsPlus(words []string) string {
	result := ""
	for i, word := range words {
		if i > 0 {
			result += ","
		}
		result += word
	}
	return result
}

func joinWordsBuilder(words []string) string {
	var builder strings.Builder
	for i, word := range words {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(word)
	}
	return builder.String()
}

// 基准测试函数（以 Benchmark 开头）
func BenchmarkJoinWordsPlus(b *testing.B) {
	words := []string{"go", "benchmark", "builder", "comparison"}
	for i := 0; i < b.N; i++ {  // b.N 由测试框架自动调整
		_ = joinWordsPlus(words)
	}
}

func BenchmarkJoinWordsBuilder(b *testing.B) {
	words := []string{"go", "benchmark", "builder", "comparison"}
	for i := 0; i < b.N; i++ {
		_ = joinWordsBuilder(words)
	}
}
```

运行基准测试：
```bash
$ go test -bench=.
goos: darwin
goarch: amd64
BenchmarkJoinWordsPlus-8          1000000    1234 ns/op
BenchmarkJoinWordsBuilder-8       5000000     234 ns/op
```

**关键点**：

- `b.N` 会由测试框架自动调整，确保结果稳定
- 比较 **ns/op**（每操作纳秒数），越小越好
- `strings.Builder` 比 `+` 快约 5 倍（避免重复分配）

### 概念 3：模糊测试（Fuzzing）

Go 1.18+ 内置模糊测试支持：

```go
// 被测函数：标准化字符串（用于 URL slug）
func normalizeSlug(input string) string {
	var builder strings.Builder
	lastWasDash := false

	for _, r := range strings.ToLower(strings.TrimSpace(input)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(r)
			lastWasDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if builder.Len() > 0 && !lastWasDash {
				builder.WriteRune('-')
				lastWasDash = true
			}
		default:
			if builder.Len() > 0 && !lastWasDash {
				builder.WriteRune('-')
				lastWasDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}

// 模糊测试函数（以 Fuzz 开头）
func FuzzNormalizeSlug(f *testing.F) {
	// 添加种子语料
	f.Add("Go Fuzzing")
	f.Add("Fuzz__Case###")
	f.Add("  Spaces  Around  ")

	// 开始模糊测试
	f.Fuzz(func(t *testing.T, input string) {
		normalized := normalizeSlug(input)
		
		// 断言：结果不能包含空格
		if strings.Contains(normalized, " ") {
			t.Fatalf("result should not contain spaces: %q", normalized)
		}
		
		// 断言：不能有连续多个 dash
		if strings.Contains(normalized, "--") {
			t.Fatalf("result should not contain consecutive dashes: %q", normalized)
		}
		
		// 断言：不能以 dash 开头或结尾
		if strings.HasPrefix(normalized, "-") || strings.HasSuffix(normalized, "-") {
			t.Fatalf("result should not start or end with dash: %q", normalized)
		}
	})
}
```

运行模糊测试：
```bash
$ go test -fuzz=FuzzNormalizeSlug
fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 3s, gathering baseline coverage: 3/3 completed
fuzz: minimizing 45-byte failing input...
```

**模糊测试的价值**：

- 自动发现你意想不到的输入组合
- 找出边界情况和特殊字符问题
- 持续运行，越久发现的问题越多

### 概念 4：测试文件组织

Go 测试文件遵循严格约定：

```
project/
├── user_service.go      # 源代码
├── user_service_test.go # 测试文件（同名 + _test.go）
└── user_service_internal_test.go  # 内部测试（访问 private）
```

**测试文件命名规则**：

| 文件后缀 | 用途 | 访问权限 |
|----------|------|----------|
| `_test.go` | 公共测试 | 只能访问导出（exported）符号 |
| `_internal_test.go` | 内部测试 | 可访问包内所有符号（包括 private） |

**测试函数类型**：

```go
// 1. 单元测试（以 Test 开头）
func TestSomething(t *testing.T) { ... }

// 2. 基准测试（以 Benchmark 开头）
func BenchmarkSomething(b *testing.B) { ... }

// 3. 模糊测试（以 Fuzz 开头）
func FuzzSomething(f *testing.F) { ... }

// 4. 示例测试（以 Example 开头，可作为文档）
func ExampleSomething() {
	// 输出会被 godoc 收录
}
```

### 概念 5：测试覆盖率（Test Coverage）

Go 内置覆盖率统计：

```bash
# 查看覆盖率
$ go test -cover
coverage: 85.7% of statements

# 生成详细报告
$ go test -coverprofile=coverage.out
$ go tool cover -html=coverage.out

# 浏览器打开查看哪些代码没被测试
$ open cover.html
```

**覆盖率不是银弹**：

- ✅ 100% 覆盖率 ≠ 没有 bug
- ✅ 低覆盖率（<50%） = 高风险
- ✅ 关注**关键路径覆盖率**，而非数字本身
- ❌ 不要为了覆盖率写无意义的测试

## 常见错误

### 错误 1：测试中包含随机性或外部依赖

```go
// ❌ 错误示例
func TestUserProfile(t *testing.T) {
	// 使用当前时间，每次运行结果不同
	user := NewUser(time.Now())
	if user.CreatedAt != time.Now() {  // 测试不稳定
		t.Fail()
	}
}

// ✅ 正确示例
func TestUserProfile(t *testing.T) {
	// 使用固定时间，测试可重现
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	user := NewUser(fixedTime)
	want := fixedTime
	if user.CreatedAt != want {
		t.Fatalf("want %v, got %v", want, user.CreatedAt)
	}
}
```

### 错误 2：一个测试做太多事情

```go
// ❌ 错误示例
func TestUserService(t *testing.T) {
	// 又测创建，又测更新，又测删除...
	user := CreateUser()
	UpdateUser(user)
	DeleteUser(user)
	// 哪个失败了？不知道
}

// ✅ 正确示例
func TestUserService_Create(t *testing.T) {
	// 只测创建
	user := CreateUser()
	// 验证创建成功
}

func TestUserService_Update(t *testing.T) {
	// 只测更新
	// ...
}

func TestUserService_Delete(t *testing.T) {
	// 只测删除
	// ...
}
```

### 错误 3：基准测试中做了多余的事情

```go
// ❌ 错误示例
func BenchmarkWrong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 在循环内分配数据，干扰性能测量
		data := loadDataFromDisk()  // 慢且不稳定
		process(data)
	}
}

// ✅ 正确示例
func BenchmarkRight(b *testing.B) {
	// 在循环外准备数据
	data := generateTestData()
	b.ResetTimer()  // 重置计时器，不包含准备时间
	for i := 0; i < b.N; i++ {
		process(data)
	}
}
```

## 动手练习

### 练习 1：完善表驱动测试

为 `normalizeSlug` 函数添加更多测试用例，覆盖以下场景：

- 空字符串输入
- 只有特殊字符的输入
- 首尾有空格的输入
- 混合大小写的输入

<details>
<summary>参考答案</summary>

```go
func TestNormalizeSlug(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "empty string", input: "", want: ""},
		{name: "only special chars", input: "!!!@@@###", want: ""},
		{name: "trim spaces", input: "  hello  ", want: "hello"},
		{name: "lowercase", input: "HELLO World", want: "hello-world"},
		{name: "mixed", input: "Go-Fuzz_Test", want: "go-fuzz-test"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeSlug(tc.input)
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}
```

</details>

### 练习 2：编写基准测试对比实现

为以下两个函数编写基准测试，比较性能差异：

```go
func plusConcat(items []string) string {
	result := ""
	for _, item := range items {
		result += item
	}
	return result
}

func builderConcat(items []string) string {
	var builder strings.Builder
	for _, item := range items {
		builder.WriteString(item)
	}
	return builder.String()
}
```

<details>
<summary>参考答案</summary>

```go
func BenchmarkPlusConcat(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e"}
	for i := 0; i < b.N; i++ {
		_ = plusConcat(items)
	}
}

func BenchmarkBuilderConcat(b *testing.B) {
	items := []string{"a", "b", "c", "d", "e"}
	for i := 0; i < b.N; i++ {
		_ = builderConcat(items)
	}
}
```

</details>

### 练习 3：编写模糊测试

为 `normalizeSlug` 编写模糊测试，确保输出符合预期规范。

<details>
<summary>参考答案</summary>

```go
func FuzzNormalizeSlug(f *testing.F) {
	f.Add("test input")
	f.Add("Go-Lang!!!")
	f.Add("  spaces  ")

	f.Fuzz(func(t *testing.T, input string) {
		result := normalizeSlug(input)
		
		// 不能有连续 dash
		if strings.Contains(result, "--") {
			t.Fatalf("consecutive dashes: %q", result)
		}
		
		// 不能有空格
		if strings.Contains(result, " ") {
			t.Fatalf("contains space: %q", result)
		}
	})
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么我的基准测试结果不稳定？

**答**：可能原因：

- 在循环内分配或准备数据（应该在循环外）
- 依赖外部资源（网络、磁盘、数据库）
- CPU 频率变化（笔记本省电模式）
- 其他进程干扰

**解决方法**：

```go
func BenchmarkStable(b *testing.B) {
	// 1. 准备工作放循环外
	data := prepareData()
	
	// 2. 重置计时器，排除准备时间
	b.ResetTimer()
	
	// 3. 多次运行取平均
	for i := 0; i < b.N; i++ {
		process(data)
	}
}
```

运行多次：`go test -bench=. -count=5`

### Q2: 表驱动测试中如何测试错误情况？

**答**：为用例添加期望错误字段：

```go
cases := []struct {
	input    string
	want     string
	wantErr  bool
	errMsg   string  // 可选：期望的错误消息
}{
	{input: "valid", want: "ok", wantErr: false},
	{input: "", want: "", wantErr: true, errMsg: "empty input"},
	{input: "invalid", want: "", wantErr: true},
}

for _, tc := range cases {
	t.Run(tc.input, func(t *testing.T) {
		got, err := process(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if tc.errMsg != "" && err.Error() != tc.errMsg {
				t.Errorf("want error %q, got %q", tc.errMsg, err.Error())
			}
		} else {
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		}
	})
}
```

### Q3: 模糊测试太慢了怎么办？

**答**：

- **设置超时**：`go test -fuzz=FuzzX -fuzztime=10s`
- **降低复杂度**：简化模糊测试函数内部的逻辑
- **减少种子**：只保留关键的种子语料
- **并行运行**：`go test -fuzz=FuzzX -parallel=4`

## 知识扩展 (选学)

### 扩展 1：表格驱动测试的高级用法

使用 `t.Run()` 创建子测试：

```go
func TestComplex(t *testing.T) {
	cases := [...]struct{
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()  // 子测试并行执行
			// ...
		})
	}
}
```

### 扩展 2：测试辅助函数

创建 `test helpers` 减少重复：

```go
// testhelper.go
func assertEqual(t *testing.T, got, want any) {
	t.Helper()  // 标记为辅助函数，错误指向调用点
	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func mustParse(t *testing.T, s string) time.Time {
	t.Helper()
	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}
```

### 扩展 3：Mock 和依赖注入

测试外部依赖（数据库、API）时使用 mock：

```go
type Database interface {
	GetUser(id int) (*User, error)
}

// 测试中使用 mock 实现
type mockDB struct{}

func (m *mockDB) GetUser(id int) (*User, error) {
	if id == 1 {
		return &User{ID: 1, Name: "test"}, nil
	}
	return nil, errors.New("not found")
}
```

## 工业界应用

### 场景：电商订单系统测试

某电商公司的订单处理系统需要高可靠性。每次修改都可能导致：

- 价格计算错误
- 库存不同步
- 优惠券滥用

**测试策略**：

```go
// 1. 表驱动测试覆盖所有折扣场景
func TestCalculateDiscount(t *testing.T) {
	cases := []struct {
		scenario     string
		orderAmount  float64
		userLevel    int
		couponCode   string
		wantDiscount float64
	}{
		{"VIP + coupon", 1000, 3, "SAVE20", 300},
		{"Normal user", 1000, 1, "", 0},
		{"Expired coupon", 1000, 2, "EXPIRED", 0},
		// 50+ 用例覆盖所有业务规则
	}
	
	for _, tc := range cases {
		t.Run(tc.scenario, func(t *testing.T) {
			got := calculateDiscount(tc.orderAmount, tc.userLevel, tc.couponCode)
			assertEqual(t, got, tc.wantDiscount)
		})
	}
}

// 2. 基准测试确保性能
func BenchmarkCalculateDiscount(b *testing.B) {
	// 双十一高并发场景
	orders := generateBenchmarkOrders()
	for i := 0; i < b.N; i++ {
		calculateDiscount(orders[i].Amount, orders[i].Level, orders[i].Coupon)
	}
}

// 3. 模糊测试发现边界情况
func FuzzValidateOrder(f *testing.F) {
	f.Add(100.50, "NYC", "12345")
	f.Fuzz(func(t *testing.T, amount float64, city string, zip string) {
		err := validateOrder(amount, city, zip)
		// 验证：负数金额应该失败
		if amount < 0 && err == nil {
			t.Fatal("negative amount should fail")
		}
	})
}
```

**效果**：

- Bug 率下降 80%
- 回归测试时间从 2 小时降至 10 分钟
- 工程师敢于重构代码

## 小结

本章介绍了 Go 测试的三大支柱：表驱动测试、基准测试、模糊测试。

### 核心概念

- **表驱动测试**：数据驱动，结构统一，易于维护
- **基准测试**：测量性能，比较实现，优化依据
- **模糊测试**：自动探索，边界发现，持续运行

### 最佳实践

1. 使用表驱动测试组织测试用例
2. 基准测试前在循环外准备数据
3. 模糊测试设置合理超时和断言
4. 测试函数小而专注，一个测试一个场景
5. 使用 `t.Helper()` 标记辅助函数

### 下一步

- 学习 `testify` 等测试库的断言功能
- 研究 Mock 模式和依赖注入
- 实践测试驱动开发（TDD）

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 表驱动测试 | Table-Driven Test | 用数据切片驱动的测试模式 |
| 基准测试 | Benchmark | 测量函数性能的测试 |
| 模糊测试 | Fuzzing | 用随机输入探索边界的测试 |
| 测试覆盖率 | Test Coverage | 被测试覆盖的代码比例 |
| 种子语料 | Seed Corpus | 模糊测试的初始输入集合 |
| 并行测试 | Parallel Test | 同时运行多个测试加快执行 |
| 辅助函数 | Helper Function | 使用 t.Helper() 标记的测试辅助代码 |
| 回归测试 | Regression Test | 确保旧功能不被破坏的测试 |
| 单元测试 | Unit Test | 测试单个函数或方法的测试 |
| Mock 测试 | Mock Test | 用模拟对象替代真实依赖的测试 |

## 源码

完整示例代码位于：[internal/advance/testing/testing.go](../../internal/advance/testing/testing.go)
