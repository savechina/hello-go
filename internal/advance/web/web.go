package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("advance", "web", Run)
}

// Run prints runnable examples for net/http handlers, middleware chains, and JSON APIs.
func Run() {
	examples := []string{
		basicHandlerExample(),
		requestResponseExample(),
		middlewareChainExample(),
		jsonAPIExample(),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

type middleware func(http.Handler) http.Handler

type apiMessage struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

func basicHandlerExample() string {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/hello", nil)

	helloHandler(recorder, request)

	return fmt.Sprintf("示例1 基础 HTTP 处理: status=%d body=%s", recorder.Code, strings.TrimSpace(recorder.Body.String()))
}

func requestResponseExample() string {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/greet?name=Gopher", nil)

	greetHandler(recorder, request)

	return fmt.Sprintf(
		"示例2 请求与响应: status=%d header=%s body=%s",
		recorder.Code,
		recorder.Header().Get("X-Request-Path"),
		strings.TrimSpace(recorder.Body.String()),
	)
}

func middlewareChainExample() string {
	handler := chainMiddlewares(
		http.HandlerFunc(greetHandler),
		requestIDMiddleware,
		authHeaderMiddleware("X-Auth-Token", "secret-token"),
	)

	deniedRecorder := httptest.NewRecorder()
	handler.ServeHTTP(deniedRecorder, httptest.NewRequest(http.MethodGet, "/greet?name=Guest", nil))

	allowedRecorder := httptest.NewRecorder()
	allowedRequest := httptest.NewRequest(http.MethodGet, "/greet?name=Middleware", nil)
	allowedRequest.Header.Set("X-Auth-Token", "secret-token")
	handler.ServeHTTP(allowedRecorder, allowedRequest)

	return fmt.Sprintf(
		"示例3 中间件链: denied=%d allowed=%d request_id=%s",
		deniedRecorder.Code,
		allowedRecorder.Code,
		allowedRecorder.Header().Get("X-Request-ID"),
	)
}

func jsonAPIExample() string {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/message?name=API", nil)

	messageAPIHandler(recorder, request)

	return fmt.Sprintf("示例4 JSON API: status=%d body=%s", recorder.Code, strings.TrimSpace(recorder.Body.String()))
}

func helloHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprint(writer, "hello, net/http learner")
}

func greetHandler(writer http.ResponseWriter, request *http.Request) {
	name := request.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		http.Error(writer, "missing name", http.StatusBadRequest)
		return
	}

	writer.Header().Set("X-Request-Path", request.URL.Path)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintf(writer, "hello, %s", name)
}

func authHeaderMiddleware(headerName string, expected string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.Header.Get(headerName) != expected {
				http.Error(writer, "unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestID := request.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = "generated-request-id"
		}

		writer.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(writer, request)
	})
}

func chainMiddlewares(handler http.Handler, middlewares ...middleware) http.Handler {
	wrapped := handler
	for index := len(middlewares) - 1; index >= 0; index-- {
		wrapped = middlewares[index](wrapped)
	}

	return wrapped
}

func messageAPIHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := request.URL.Query().Get("name")
	if strings.TrimSpace(name) == "" {
		name = "gopher"
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(apiMessage{
		Message: "hello, " + name,
		Path:    request.URL.Path,
	})
}
