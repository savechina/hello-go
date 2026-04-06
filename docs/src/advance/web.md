# Web 开发（Web Development）

## 开篇故事

想象你开了一家餐厅。客人进门（发起请求），服务员接待（Handler），点菜（读取请求参数），厨房做菜（业务逻辑），上菜（返回响应）。这个流程简单直观，但要做好却需要很多细节：

- 服务员要听懂客人的要求（解析请求）
- 厨房要按标准做菜（业务逻辑）
- 上菜前要摆盘（设置响应头）
- 遇到投诉要处理（错误处理）

Go 语言的 `net/http` 标准库就像这套餐厅运营系统。它设计简洁，但功能完整：Handler 是服务员，Request 是客人点单，Response 是端上去的菜，Middleware 是经理（可以在服务员和客人之间做额外处理）。

很多初学者一上来就学框架（Gin、Echo），但框架的本质是对标准库的封装。理解标准库，就像学会了餐厅运营的基本功，换到任何框架都能快速上手。不理解标准库，就像只学过某个连锁店的点餐系统，换个店就不会工作了。

本章从 `net/http` 出发，带你理解 Web 服务的核心概念：Handler、Request、Response、Middleware。学完这些，你再去看任何 Web 框架，都会发现"原来如此"。

## 本章适合谁

- ✅ 已掌握 Go 基础语法（函数、结构体、接口）的开发者
- ✅ 想理解 Web 服务工作原理的学习者
- ✅ 准备学习 Web 框架但想先打好基础的工程师
- ✅ 需要编写 HTTP API 或 Web 服务的技术人员

如果你还没有写过基本的 Go 程序，建议先完成基础章节。

## 你会学到什么

学完本章后，你将能够：

1. **编写 HTTP Handler**：理解 `http.Handler` 接口的核心职责
2. **处理请求和响应**：读取查询参数、设置响应头、返回不同格式
3. **实现中间件**：编写日志、鉴权、请求追踪等横切逻辑
4. **构建 JSON API**：序列化数据、设置 Content-Type、处理错误
5. **测试 HTTP 代码**：使用 `httptest` 无需启动端口测试 Handler

## 前置要求

在开始本章之前，请确保你已经掌握：

- Go 基础语法（函数、结构体、接口）
- 错误处理基础
- JSON 基础（encoding/json）
- 对 HTTP 协议有基本概念（请求、响应、状态码）

## 第一个例子

让我们从最简单的 Handler 开始：

```go
package main

import (
	"fmt"
	"net/http"
)

// 最简单的 Handler：实现 ServeHTTP 方法
type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. 设置响应头（必须在 Write 之前）
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	
	// 2. 写入响应体
	fmt.Fprint(w, "hello, net/http learner")
}

// 也可以用函数实现 Handler
func helloFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "hello from function handler")
}

func main() {
	// 注册路由
	http.Handle("/hello", helloHandler{})
	http.HandleFunc("/hello-func", helloFunc)
	
	// 启动服务器
	fmt.Println("Server starting at :8080")
	http.ListenAndServe(":8080", nil)
}
```

测试：
```bash
$ curl http://localhost:8080/hello
hello, net/http learner

$ curl http://localhost:8080/hello-func
hello from function handler
```

这个例子展示了 Handler 的核心：**接收请求，返回响应**。`http.Handler` 接口只有一个方法：

```go
type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}
```

## 原理解析

### 概念 1：http.Handler 接口

`http.Handler` 是 Go Web 的基石：

```go
type Handler interface {
	ServeHTTP(w ResponseWriter, r *Request)
}
```

**两种实现方式**：

```go
// 方式 1：结构体实现（适合有状态的 Handler）
type greetHandler struct {
	prefix string
}

func (h greetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s, %s", h.prefix, r.URL.Query().Get("name"))
}

// 方式 2：函数实现（适合简单无状态逻辑）
func greetFunc(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	fmt.Fprintf(w, "hello, %s", name)
}

// http.HandlerFunc 是适配器类型
type HandlerFunc func(ResponseWriter, *Request)

func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
	f(w, r)  // 调用函数本身
}
```

**注册路由**：

```go
// 注册 Handler
http.Handle("/hello", helloHandler{})

// 注册 HandlerFunc（自动适配）
http.HandleFunc("/greet", greetFunc)

// 使用匿名函数
http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "pong")
})
```

### 概念 2：请求处理（Request Handling）

`*http.Request` 包含所有请求信息：

```go
func detailedHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 读取 URL 路径
	path := r.URL.Path  // "/api/users"
	
	// 2. 读取查询参数
	name := r.URL.Query().Get("name")
	page := r.URL.Query().Get("page")
	
	// 3. 读取请求头
	userAgent := r.Header.Get("User-Agent")
	authToken := r.Header.Get("Authorization")
	
	// 4. 读取请求方法
	method := r.Method  // "GET", "POST", "PUT", "DELETE"
	
	// 5. 读取请求体（POST/PUT）
	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
		defer r.Body.Close()
	}
	
	// 6. 其他信息
	remoteAddr := r.RemoteAddr  // 客户端地址
	tls := r.TLS != nil         // 是否 HTTPS
	
	fmt.Fprintf(w, "path=%s, name=%s, method=%s", path, name, method)
}
```

**查询参数处理**：

```go
// URL: /search?q=golang&page=2&limit=10
query := r.URL.Query()
q := query.Get("q")      // "golang"
page := query.Get("page") // "2"
// 不存在的参数返回空字符串
missing := query.Get("missing") // ""

// 多值参数
// URL: /tags?tag=go&tag=rust
tags := query["tag"]  // []string{"go", "rust"}
```

### 概念 3：响应处理（Response Handling）

`http.ResponseWriter` 用于构建响应：

```go
func responseHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 设置响应头（必须在 WriteHeader 或 Write 之前）
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Custom-Header", "custom-value")
	w.Header().Set("X-Request-ID", "abc123")
	
	// 2. 设置状态码（默认 200 OK）
	w.WriteHeader(http.StatusOK)  // 200
	// w.WriteHeader(http.StatusNotFound)  // 404
	// w.WriteHeader(http.StatusInternalServerError)  // 500
	
	// 3. 写入响应体
	fmt.Fprint(w, `{"status":"ok"}`)
	
	// 注意：WriteHeader 只能调用一次
	// 第一次调用后，后续的 WriteHeader 调用无效
}
```

**便捷函数**：

```go
// 快速返回错误响应
func errorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "something went wrong", http.StatusInternalServerError)
	// 等价于：
	// w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// w.WriteHeader(500)
	// fmt.Fprint(w, "something went wrong")
}

// 重定向
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/new-location", http.StatusMovedPermanently)
}
```

### 概念 4：中间件（Middleware）

中间件是"包装 Handler 的 Handler"：

```go
// 中间件类型
type middleware func(http.Handler) http.Handler

// 日志中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 请求前处理
		start := time.Now()
		fmt.Printf("[%s] %s %s\n", start.Format(time.RFC3339), r.Method, r.URL.Path)
		
		// 调用下一个 Handler
		next.ServeHTTP(w, r)
		
		// 请求后处理
		fmt.Printf("[%s] completed in %v\n", time.Now(), time.Since(start))
	})
}

// 鉴权中间件
func authHeaderMiddleware(headerName string, expected string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(headerName)
			if token != expected {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// 请求 ID 中间件
func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateID()  // 生成唯一 ID
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}
```

**中间件链**：

```go
// 手动嵌套
handler := loggingMiddleware(
	authHeaderMiddleware("X-Auth", "secret")(
		requestIDMiddleware(http.HandlerFunc(handlerFunc)),
	)
)

// 使用辅助函数
func chainMiddlewares(handler http.Handler, middlewares ...middleware) http.Handler {
	wrapped := handler
	// 从后往前包装（类似洋葱模型）
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}

// 使用
handler := chainMiddlewares(
	http.HandlerFunc(handlerFunc),
	loggingMiddleware,
	requestIDMiddleware,
	authHeaderMiddleware("X-Auth", "secret"),
)
```

### 概念 5：JSON API

现代 Web 服务最常用的响应格式：

```go
type apiMessage struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

func messageAPIHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 验证请求方法
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// 2. 读取参数
	name := r.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		name = "gopher"  // 默认值
	}
	
	// 3. 构建响应数据
	response := apiMessage{
		Message: "hello, " + name,
		Path:    r.URL.Path,
	}
	
	// 4. 设置 JSON Content-Type
	w.Header().Set("Content-Type", "application/json")
	
	// 5. 编码并返回
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
```

**错误响应的 JSON 格式**：

```go
type errorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details,omitempty"`
}

func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{
		Error: message,
		Code:  statusCode,
	})
}

// 使用
if err != nil {
	writeJSONError(w, http.StatusBadRequest, "invalid input")
	return
}
```

### 概念 6：httptest 测试

无需启动端口即可测试 Handler：

```go
func TestHelloHandler(t *testing.T) {
	// 1. 创建响应记录器
	recorder := httptest.NewRecorder()
	
	// 2. 创建测试请求
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)
	
	// 3. 调用 Handler
	helloHandler{}.ServeHTTP(recorder, request)
	
	// 4. 验证响应
	if recorder.Code != http.StatusOK {
		t.Errorf("want status %d, got %d", http.StatusOK, recorder.Code)
	}
	
	if recorder.Body.String() != "hello, net/http learner" {
		t.Errorf("want body %q, got %q", "hello, net/http learner", recorder.Body.String())
	}
	
	// 5. 验证响应头
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("want content-type %q, got %q", "text/plain; charset=utf-8", contentType)
	}
}
```

**测试带参数的请求**：

```go
func TestGreetHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/greet?name=Gopher", nil)
	
	greetHandler(recorder, request)
	
	if recorder.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Code)
	}
	
	if !strings.Contains(recorder.Body.String(), "Gopher") {
		t.Errorf("expected greeting for Gopher, got: %s", recorder.Body.String())
	}
}
```

**测试中间件**：

```go
func TestAuthMiddleware(t *testing.T) {
	handler := authHeaderMiddleware("X-Auth", "secret-token")(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "success")
		}),
	)
	
	// 测试无 token 的情况
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(recorder, request)
	
	if recorder.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", recorder.Code)
	}
	
	// 测试有正确 token 的情况
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Auth", "secret-token")
	handler.ServeHTTP(recorder, request)
	
	if recorder.Code != http.StatusOK {
		t.Errorf("want 200, got %d", recorder.Code)
	}
}
```

## 常见错误

### 错误 1：在 Write 之后设置响应头

```go
// ❌ 错误示例
func wrongHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")  // 先写入
	w.Header().Set("Content-Type", "text/plain")  // 无效了！
}

// ✅ 正确示例
func rightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")  // 先设置头
	fmt.Fprint(w, "hello")
}
```

### 错误 2：忘记检查请求方法

```go
// ❌ 错误示例
func createUser(w http.ResponseWriter, r *http.Request) {
	// 没有检查 Method，GET 请求也能创建用户
	var input User
	json.NewDecoder(r.Body).Decode(&input)
	// ...
}

// ✅ 正确示例
func createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// ...
}
```

### 错误 3：忘记关闭请求体

```go
// ❌ 错误示例
func readBody(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)  // 忘记 Close，可能泄漏
	// ...
}

// ✅ 正确示例
func readBody(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()  // 确保关闭
	body, _ := io.ReadAll(r.Body)
	// ...
}
```

### 错误 4：JSON 响应忘记设置 Content-Type

```go
// ❌ 错误示例
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	// 客户端可能无法正确解析
}

// ✅ 正确示例
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
```

## 动手练习

### 练习 1：实现查询参数验证

为 greetHandler 添加验证逻辑：如果 name 参数为空，返回 400 错误。

<details>
<summary>参考答案</summary>

```go
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		http.Error(w, "missing name parameter", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "hello, %s", name)
}
```

</details>

### 练习 2：实现 JSON API Handler

编写一个 Handler，接收 GET 请求，返回 JSON 格式的问候信息。

<details>
<summary>参考答案</summary>

```go
type greetResponse struct {
	Greeting string `json:"greeting"`
	Name     string `json:"name"`
}

func jsonGreetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	
	response := greetResponse{
		Greeting: "Hello",
		Name:     name,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
```

</details>

### 练习 3：实现日志中间件

编写一个中间件，记录每个请求的方法、路径、耗时。

<details>
<summary>参考答案</summary>

```go
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// 调用下一个 Handler
		next.ServeHTTP(w, r)
		
		// 记录日志
		duration := time.Since(start)
		fmt.Printf("[%s] %s %s - %v\n", 
			start.Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			duration)
	})
}
```

</details>

## 故障排查 (FAQ)

### Q1: 为什么响应头设置没有生效？

**答**：最常见原因是设置时机太晚。响应头必须在第一次 `Write()` 或 `WriteHeader()` 之前设置：

```go
// ❌ 错误：Write 后设置头无效
fmt.Fprint(w, "body")
w.Header().Set("X-Custom", "value")  // 无效

// ✅ 正确：先设置头
w.Header().Set("X-Custom", "value")
fmt.Fprint(w, "body")
```

### Q2: 为什么 Handler 返回 404？

**答**：检查路由注册：

```go
// ❌ 错误：路径不匹配
http.Handle("/api", handler)  // 只匹配 /api
// /api/ 或 /api/users 会返回 404

// ✅ 正确
http.Handle("/api/", handler)  // 匹配 /api/ 及其子路径
```

或使用 ServeMux：

```go
mux := http.NewServeMux()
mux.HandleFunc("/api/users", handler)
```

### Q3: 如何调试 Handler 的逻辑？

**答**：使用 httptest 记录完整响应：

```go
recorder := httptest.NewRecorder()
request := httptest.NewRequest("GET", "/test?debug=1", nil)
handler.ServeHTTP(recorder, request)

// 打印所有信息
fmt.Println("Status:", recorder.Code)
fmt.Println("Headers:", recorder.Header())
fmt.Println("Body:", recorder.Body.String())
```

## 知识扩展 (选学)

### 扩展 1：使用 httprouter 或 chi

标准库的 ServeMux 功能有限，生产环境常用更强大的路由库：

```go
// chi 示例
r := chi.NewRouter()
r.Get("/users/{id}", getUserHandler)
r.Post("/users", createUserHandler)
r.Put("/users/{id}", updateUserHandler)
r.Delete("/users/{id}", deleteUserHandler)
```

### 扩展 2：Graceful Shutdown

优雅关闭服务器：

```go
server := &http.Server{Addr: ":8080", Handler: mux}

go func() {
	<-shutdownChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}()

server.ListenAndServe()
```

### 扩展 3：HTTPS 配置

```go
server := &http.Server{
	Addr:      ":443",
	Handler:   mux,
	TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
}

server.ListenAndServeTLS("cert.pem", "key.pem")
```

## 工业界应用

### 场景：微服务 API 网关

某公司的 API 网关需要处理：

- 请求鉴权（JWT 验证）
- 限流（rate limiting）
- 请求日志（审计）
- 响应缓存（性能）

**中间件架构**：

```go
// 定义 Handler 链
handler := chainMiddlewares(
	apiRouter,
	loggingMiddleware,           // 最外层：记录所有请求
	recoveryMiddleware,          // panic 恢复
	rateLimitMiddleware(100),    // 限流：100 req/s
	authMiddleware(jwtVerifier), // JWT 验证
	corsMiddleware,              // CORS 支持
	compressionMiddleware,       // Gzip 压缩
)

// 每个 Handler 专注于业务逻辑
apiRouter.HandleFunc("/users", listUsers)
apiRouter.HandleFunc("/users/{id}", getUser)
```

**效果**：

- 业务逻辑和横切关注点分离
- 中间件可复用、可测试
- 新增功能只需添加中间件

## 小结

本章介绍了 Go Web 开发的核心概念：Handler、Request、Response、Middleware。

### 核心概念

- **http.Handler**：Web 服务的基石接口
- **Request**：包含所有请求信息
- **ResponseWriter**：用于构建响应
- **Middleware**：包装 Handler 的横切逻辑
- **httptest**：无需端口的测试工具

### 最佳实践

1. 先设置响应头，再写入响应体
2. 检查请求方法再处理业务
3. JSON 响应必须设置 Content-Type
4. 使用 defer 关闭请求体
5. 用 httptest 测试而非启动真实服务器

### 下一步

- 学习 Gin 或 Echo 等框架
- 实践 RESTful API 设计
- 学习 WebSocket 实时通信

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| Handler | Handler | 处理 HTTP 请求的接口 |
| 中间件 | Middleware | 包装 Handler 的横切逻辑 |
| 请求 | Request | HTTP 请求对象 |
| 响应 | Response | HTTP 响应对象 |
| 路由 | Routing | URL 路径到 Handler 的映射 |
| 查询参数 | Query Parameter | URL ? 后的参数 |
| 请求头 | Request Header | 请求元数据 |
| 响应头 | Response Header | 响应元数据 |
| 状态码 | Status Code | HTTP 响应状态标识 |
| Content-Type | Content-Type | 响应体格式声明 |

## 源码

完整示例代码位于：[internal/advance/web/web.go](../../internal/advance/web/web.go)
