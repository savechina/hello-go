package web

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlersAndMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		handler          http.Handler
		method           string
		target           string
		headers          map[string]string
		wantStatus       int
		wantBodyContains string
		wantHeaderKey    string
		wantHeaderValue  string
	}{
		{
			name:             "basic handler",
			handler:          http.HandlerFunc(helloHandler),
			method:           http.MethodGet,
			target:           "/hello",
			wantStatus:       http.StatusOK,
			wantBodyContains: "hello, net/http learner",
		},
		{
			name:             "request response handler",
			handler:          http.HandlerFunc(greetHandler),
			method:           http.MethodGet,
			target:           "/greet?name=Gopher",
			wantStatus:       http.StatusOK,
			wantBodyContains: "hello, Gopher",
			wantHeaderKey:    "X-Request-Path",
			wantHeaderValue:  "/greet",
		},
		{
			name:             "middleware unauthorized",
			handler:          chainMiddlewares(http.HandlerFunc(greetHandler), authHeaderMiddleware("X-Auth-Token", "secret-token")),
			method:           http.MethodGet,
			target:           "/greet?name=Guest",
			wantStatus:       http.StatusUnauthorized,
			wantBodyContains: "unauthorized",
		},
		{
			name: "middleware authorized",
			handler: chainMiddlewares(
				http.HandlerFunc(greetHandler),
				requestIDMiddleware,
				authHeaderMiddleware("X-Auth-Token", "secret-token"),
			),
			method:           http.MethodGet,
			target:           "/greet?name=Middleware",
			headers:          map[string]string{"X-Auth-Token": "secret-token"},
			wantStatus:       http.StatusOK,
			wantBodyContains: "hello, Middleware",
			wantHeaderKey:    "X-Request-ID",
			wantHeaderValue:  "generated-request-id",
		},
		{
			name:             "json api handler",
			handler:          http.HandlerFunc(messageAPIHandler),
			method:           http.MethodGet,
			target:           "/api/message?name=API",
			wantStatus:       http.StatusOK,
			wantBodyContains: "\"message\":\"hello, API\"",
			wantHeaderKey:    "Content-Type",
			wantHeaderValue:  "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.target, nil)
			for key, value := range tt.headers {
				request.Header.Set(key, value)
			}

			recorder := httptest.NewRecorder()
			tt.handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, recorder.Code)
			}
			if !strings.Contains(recorder.Body.String(), tt.wantBodyContains) {
				t.Fatalf("expected body to contain %q, got %q", tt.wantBodyContains, recorder.Body.String())
			}
			if tt.wantHeaderKey != "" && recorder.Header().Get(tt.wantHeaderKey) != tt.wantHeaderValue {
				t.Fatalf("expected header %s=%q, got %q", tt.wantHeaderKey, tt.wantHeaderValue, recorder.Header().Get(tt.wantHeaderKey))
			}
		})
	}
}

func TestRunOutput(t *testing.T) {
	output := captureOutput(t, Run)

	tests := []struct {
		name string
		want string
	}{
		{name: "basic handler example", want: "示例1 基础 HTTP 处理"},
		{name: "request response example", want: "示例2 请求与响应"},
		{name: "middleware example", want: "示例3 中间件链"},
		{name: "json api example", want: "示例4 JSON API"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(output, tt.want) {
				t.Fatalf("expected output to contain %q, got %q", tt.want, output)
			}
		})
	}
}

func captureOutput(t *testing.T, runner func()) string {
	t.Helper()

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	os.Stdout = w
	runner()
	_ = w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read output: %v", err)
	}

	return buf.String()
}
