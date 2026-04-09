# Web 服务实战（Web Service）

## 开篇场景

想象你要开发一个任务管理系统的后端 API。用户可以通过手机 App 添加任务、查看任务列表、标记完成状态。这个系统需要处理多个客户端同时访问的情况，还要记录每次请求的日志，方便排查问题。

听起来很简单，但要实现一个生产级别的服务，你需要考虑这些问题：

- 多个用户同时添加任务时，数据会不会混乱？
- API 返回的数据格式是否标准（JSON）？
- 如何记录每个请求的方法和路径？
- 不启动服务器也能测试 API 是否正常工作？

本实战项目用 Go 标准库 `net/http` 实现了一个完整的 RESTful API 示例，包含线程安全的数据存储、中间件模式和测试技巧。代码不到 100 行，却涵盖了 Web 服务开发的核心技能。

## 项目概览

这个实战项目位于 `internal/awesome/webservice/webservice.go`，实现了以下功能：

1. **Task 结构体**：定义任务的数据模型，包含 ID、标题、完成状态
2. **线程安全存储**：使用 `sync.RWMutex` 保护共享数据，支持并发读写
3. **RESTful Handler**：实现 GET 列表和 POST 创建两个核心接口
4. **中间件模式**：日志中间件演示请求追踪的实现方式
5. **httptest 测试**：无需启动端口即可验证 Handler 行为

## 概念说明

### 1. RESTful API 设计原则

REST（Representational State Transfer）是一种 Web 服务架构风格。它的核心思想是用统一的接口操作资源：

| HTTP 方法 | 用途 | 示例 |
|----------|------|------|
| GET | 获取资源列表或单个资源 | GET /tasks 获取所有任务 |
| POST | 创建新资源 | POST /tasks 创建新任务 |
| PUT | 更新现有资源 | PUT /tasks/1 更新任务 1 |
| DELETE | 删除资源 | DELETE /tasks/1 删除任务 1 |

本项目实现了 GET 和 POST 两个方法，涵盖了最常见的 API 操作。

### 2. 线程安全（Thread Safety）

Web 服务天然是并发环境。多个请求可能同时到达，如果共享数据没有保护，会导致数据竞争（race condition）：

```go
// ❌ 不安全的写法
tasks := []Task{}
tasks = append(tasks, newTask)  // 并发调用会导致数据丢失或损坏

// ✅ 使用互斥锁保护
var mu sync.Mutex
mu.Lock()
tasks = append(tasks, newTask)
mu.Unlock()
```

Go 提供了两种互斥锁：
- `sync.Mutex`：读写都互斥，适合写多读少的场景
- `sync.RWMutex`：读共享、写互斥，适合读多写少的场景（本项目采用）

### 3. Handler 函数

Handler 是处理 HTTP 请求的核心组件。Go 的 `http.HandlerFunc` 类型让普通函数可以直接作为 Handler：

```go
type HandlerFunc func(ResponseWriter, *Request)
```

Handler 的职责：
1. 解析请求参数和请求体
2. 执行业务逻辑
3. 设置响应头（Content-Type、状态码）
4. 编码并返回响应数据

### 4. 中间件模式（Middleware Pattern）

中间件是"包装 Handler 的 Handler"。它可以在请求到达业务 Handler 前后添加额外处理，比如：

- 日志记录：记录请求方法、路径、耗时
- 身份认证：验证用户是否有权限
- 请求限流：防止恶意请求攻击
- 错误恢复：捕获 panic 防止服务崩溃

中间件的典型结构：

```go
func middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 请求前处理
        next(w, r)  // 调用下一个 Handler
        // 请求后处理
    }
}
```

### 5. httptest 测试技巧

测试 Web Handler 传统方式是启动服务器，用 curl 或 Postman 发送请求。这种方式有几个问题：

- 测试速度慢（每次都要启动端口）
- 难以自动化（需要外部工具）
- 无法验证内部细节（比如响应头）

`net/http/httptest` 提供了不启动端口就能测试 Handler 的能力：

- `httptest.NewRecorder()`：模拟 ResponseWriter，记录响应内容
- `httptest.NewRequest()`：创建模拟请求，设置方法和路径

## 代码示例

### 示例 1：定义数据模型

Task 结构体定义了任务的数据模型：

```go
// Task represents a todo item.
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}
```

要点解析：
- JSON 标签（`json:"id"`）指定了序列化时的字段名
- `int` 类型用于唯一标识，适合数据库存储
- `bool` 类型表示完成状态，语义清晰

### 示例 2：线程安全存储

Store 结构体使用读写锁保护任务列表：

```go
// Store holds tasks with thread-safe access.
type Store struct {
	mu     sync.RWMutex
	tasks  []Task
	nextID int
}

func (s *Store) List() []Task {
	s.mu.RLock()         // 读锁：允许多个并发读取
	defer s.mu.RUnlock() // defer 确保函数返回时释放锁
	return s.tasks
}

func (s *Store) Add(title string) Task {
	s.mu.Lock()          // 写锁：独占访问
	defer s.mu.Unlock()  // defer 确保函数返回时释放锁
	t := Task{ID: s.nextID, Title: title}
	s.nextID++
	s.tasks = append(s.tasks, t)
	return t
}
```

要点解析：
- `RLock()` 用于读操作，多个 goroutine 可以同时持有读锁
- `Lock()` 用于写操作，会阻塞所有其他锁请求
- `defer` 确保锁一定会释放，避免忘记 Unlock 导致死锁
- `nextID` 在锁保护下递增，保证 ID 唯一性

### 示例 3：Handler 实现

listHandler 处理 GET 请求，返回任务列表：

```go
listHandler := func(w http.ResponseWriter, r *http.Request) {
	tasks := store.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
```

addTask 处理 POST 请求，创建新任务：

```go
addTask := func(w http.ResponseWriter, r *http.Request) {
	var req struct{ Title string }
	json.NewDecoder(r.Body).Decode(&req)
	t := store.Add(req.Title)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}
```

要点解析：
- `json.NewEncoder(w).Encode()`：直接向 ResponseWriter 编码 JSON
- `json.NewDecoder(r.Body).Decode()`：从请求体解码 JSON
- `w.WriteHeader(http.StatusCreated)`：设置 201 状态码（资源创建成功）
- 响应头必须在 WriteHeader 或 Encode 之前设置

### 示例 4：中间件实现

日志中间件记录每个请求：

```go
loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("  [LOG] %s %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}
```

要点解析：
- 接收 `http.HandlerFunc` 作为参数，返回一个新的 `http.HandlerFunc`
- 可以在 `next(w, r)` 前后添加处理逻辑
- 中间件可以叠加使用，形成处理链

### 示例 5：httptest 测试

使用 httptest 测试 Handler：

```go
rec := httptest.NewRecorder()
req := httptest.NewRequest("GET", "/tasks", nil)
loggingMiddleware(listHandler)(rec, req)
fmt.Printf("    响应状态: %d\n", rec.Code)
```

要点解析：
- `NewRecorder()` 创建响应记录器，可以读取响应内容
- `NewRequest()` 创建模拟请求，指定方法和路径
- Handler 直接调用，无需启动服务器
- `rec.Code` 获取响应状态码，`rec.Body.String()` 获取响应体

## 知识点总结

### 核心技能

| 技能 | 本项目体现 | 实际应用 |
|------|----------|---------|
| 结构体定义 | Task 结构体 | 定义 API 数据模型 |
| 互斥锁使用 | sync.RWMutex | 保护并发访问的共享数据 |
| JSON 序列化 | json.Encoder/Decoder | RESTful API 数据交换格式 |
| Handler 编写 | listHandler/addTask | 处理 HTTP 请求的核心逻辑 |
| 中间件模式 | loggingMiddleware | 请求日志、鉴权、限流等横切逻辑 |
| httptest 测试 | NewRecorder/NewRequest | 单元测试和集成测试 |

### 最佳实践

1. **锁的范围最小化**：只在访问共享数据时持有锁，避免阻塞其他操作
2. **defer 释放锁**：确保锁一定会释放，即使发生 panic
3. **先设置响应头**：Content-Type 必须在写入响应体之前设置
4. **使用状态码语义**：200 成功、201 创建、400 客户端错误、500 服务端错误
5. **中间件函数签名**：统一使用 `func(http.HandlerFunc) http.HandlerFunc`

### 常见陷阱

| 陷阱 | 表现 | 解决方法 |
|------|------|---------|
| 锁忘记释放 | 死锁，程序卡住 | 使用 defer Unlock |
| 响应头设置时机错误 | Content-Type 不生效 | 在 WriteHeader/Encode 前设置 |
| 未关闭请求体 | 内存泄漏 | 使用 defer r.Body.Close() |
| 未检查请求方法 | GET 能触发 POST 操作 | 先判断 r.Method |

## 练习题与思考题

### 练习 1：添加完成状态更新接口

在 Store 中添加 `Complete(id int)` 方法，实现标记任务完成的功能。Handler 应该用什么 HTTP 方法？

<details>
<summary>参考答案</summary>

```go
func (s *Store) Complete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks[i].Completed = true
			return true
		}
	}
	return false
}

// Handler: PUT /tasks/{id}
completeHandler := func(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)  // 从路径提取 ID
	if !store.Complete(id) {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
```

应该用 PUT 方法，因为这是更新现有资源。
</details>

### 练习 2：实现超时中间件

编写一个中间件，如果 Handler 执行超过 1 秒，返回 504 Gateway Timeout。

<details>
<summary>参考答案</summary>

```go
func timeoutMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		
		// 用带超时的 context 替换原请求
		r = r.WithContext(ctx)
		
		done := make(chan struct{})
		go func() {
			next(w, r)
			close(done)
		}()
		
		select {
		case <-done:
			// Handler 正常完成
		case <-ctx.Done():
			http.Error(w, "request timeout", http.StatusGatewayTimeout)
		}
	}
}
```

注意：这个实现使用了 context 和 goroutine，是更高级的模式。
</details>

### 练习 3：编写 httptest 单元测试

为 addTask Handler 编写完整的单元测试，验证：
- 状态码是 201
- 响应体包含新创建的任务
- Content-Type 是 application/json

<details>
<summary>参考答案</summary>

```go
func TestAddTask(t *testing.T) {
	store := &Store{}
	
	rec := httptest.NewRecorder()
	body := strings.NewReader(`{"title":"Test Task"}`)
	req := httptest.NewRequest("POST", "/tasks", body)
	
	addTask(rec, req)
	
	// 验证状态码
	if rec.Code != http.StatusCreated {
		t.Errorf("want status 201, got %d", rec.Code)
	}
	
	// 验证 Content-Type
	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("want application/json, got %s", ct)
	}
	
	// 验证响应体
	var task Task
	if err := json.Unmarshal(rec.Body.Bytes(), &task); err != nil {
		t.Fatalf("parse response: %v", err)
	}
	if task.Title != "Test Task" {
		t.Errorf("want title 'Test Task', got '%s'", task.Title)
	}
}
```

</details>

### 思考题 1：为什么用 RWMutex 而不是 Mutex？

分析这个场景：任务管理系统的读操作（查看列表）频率远高于写操作（添加任务）。如果用普通 Mutex，读请求之间也会相互阻塞，性能会下降。RWMutex 的读锁可以共享，多个用户同时查看任务列表不会互相等待。

### 思考题 2：中间件可以叠加吗？顺序有什么影响？

中间件可以无限叠加，形成"洋葱模型"：

```go
handler := loggingMiddleware(
	authMiddleware(
		timeoutMiddleware(realHandler),
	),
)
```

请求执行顺序：logging → auth → timeout → realHandler → timeout → auth → logging

顺序很重要：
- 认证失败应该在超时判断之前（避免浪费超时检测）
- 日志应该在最外层（记录所有请求，包括被拒绝的）

### 思考题 3：如何防止恶意请求发送超大 JSON？

攻击者可能发送超大 JSON 请求体，消耗服务器内存。防护措施：

1. 限制请求体大小：

```go
r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)  // 最大 1MB
```

2. 在解码前检查 Content-Length：

```go
if r.ContentLength > 1024*1024 {
	http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
	return
}
```

## 源码位置

完整代码位于：[internal/awesome/webservice/webservice.go](../../internal/awesome/webservice/webservice.go)

运行示例：

```bash
go run cmd/hello/main.go awesome webservice
```

## 扩展阅读

如果想深入学习 Web 服务开发，可以继续探索：

- **路由进阶**：使用 `gorilla/mux` 或 `chi` 实现更复杂的路由规则
- **认证鉴权**：JWT（JSON Web Token）实现无状态认证
- **数据库集成**：将 Store 替换为真实数据库（SQLite、PostgreSQL）
- **Graceful Shutdown**：优雅关闭服务器，不中断正在处理的请求
- **Rate Limiting**：使用令牌桶算法限制请求频率

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 线程安全 | Thread Safety | 多个 goroutine 同时访问不会导致数据错误 |
| 互斥锁 | Mutex | 保证同一时刻只有一个 goroutine 访问资源 |
| 读写锁 | RWMutex | 读操作可共享，写操作独占 |
| Handler | Handler | 处理 HTTP 请求的函数 |
| 中间件 | Middleware | 包装 Handler 的横切逻辑组件 |
| RESTful | RESTful | 基于 HTTP 方语法的资源操作风格 |
| 序列化 | Serialization | 将数据结构转换为 JSON 等格式 |
| httptest | httptest | Go 标准库的 HTTP 测试工具 |