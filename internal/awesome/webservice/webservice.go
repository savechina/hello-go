package webservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"

	"hello/internal/chapters"
)

// Task represents a todo item.
type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// Store holds tasks with thread-safe access.
type Store struct {
	mu     sync.RWMutex
	tasks  []Task
	nextID int
}

func (s *Store) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tasks
}

func (s *Store) Add(title string) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	t := Task{ID: s.nextID, Title: title}
	s.nextID++
	s.tasks = append(s.tasks, t)
	return t
}

// Run demonstrates a RESTful API server with handlers and middleware.
func Run() {
	fmt.Println("=== 实战项目：Web 服务 (Web Service) ===")

	store := &Store{}

	// 示例1: Handler 函数 (Handler function)
	listHandler := func(w http.ResponseWriter, r *http.Request) {
		tasks := store.List()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}

	// 示例2: 添加任务 (Add task)
	addTask := func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Title string }
		json.NewDecoder(r.Body).Decode(&req)
		t := store.Add(req.Title)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	}

	// 示例3: Middleware 中间链 (Middleware chain)
	loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("  [LOG] %s %s\n", r.Method, r.URL.Path)
			next(w, r)
		}
	}

	// 演示：模拟请求
	fmt.Println("  示例1: GET /tasks")
	store.Add("Learn Go")
	store.Add("Build API")
	tasks := store.List()
	fmt.Printf("    返回: %+v\n", tasks)

	fmt.Println("  示例2: POST /tasks {title: \"Test middleware\"}")
	newTask := store.Add("Test middleware")
	fmt.Printf("    返回: %+v\n", newTask)

	fmt.Println("  示例3: Middleware 日志 + httptest")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/tasks", nil)
	loggingMiddleware(listHandler)(rec, req)
	fmt.Printf("    响应状态: %d\n", rec.Code)

	// Use addTask to avoid unused warning
	_ = addTask
}

func init() {
	chapters.Register("awesome", "webservice", Run)
}
