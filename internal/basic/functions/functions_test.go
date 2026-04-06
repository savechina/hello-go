package functions

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFunctionHelpers(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func(t *testing.T)
	}{
		{
			name: "greet",
			checkFunc: func(t *testing.T) {
				if got := greet("Go"); got != "Hello, Go" {
					t.Fatalf("expected greeting %q, got %q", "Hello, Go", got)
				}
			},
		},
		{
			name: "divide and remainder",
			checkFunc: func(t *testing.T) {
				quotient, remainder, err := divideAndRemainder(17, 5)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if quotient != 3 || remainder != 2 {
					t.Fatalf("expected 3 and 2, got %d and %d", quotient, remainder)
				}
			},
		},
		{
			name: "divide by zero",
			checkFunc: func(t *testing.T) {
				_, _, err := divideAndRemainder(10, 0)
				if err == nil {
					t.Fatal("expected error for divide by zero")
				}
			},
		},
		{
			name: "variadic sum",
			checkFunc: func(t *testing.T) {
				if got := sumAll(1, 2, 3, 4); got != 10 {
					t.Fatalf("expected sum 10, got %d", got)
				}
			},
		},
		{
			name: "named return metrics",
			checkFunc: func(t *testing.T) {
				area, perimeter := rectangleMetrics(3, 4)
				if area != 12 || perimeter != 14 {
					t.Fatalf("expected area 12 and perimeter 14, got %.1f and %.1f", area, perimeter)
				}
			},
		},
		{
			name: "closure counter",
			checkFunc: func(t *testing.T) {
				counter := makeCounter(5)
				if first, second := counter(), counter(); first != 6 || second != 7 {
					t.Fatalf("expected counter sequence 6,7 got %d,%d", first, second)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.checkFunc)
	}
}

func TestRunOutput(t *testing.T) {
	output := captureOutput(t, Run)

	tests := []struct {
		name string
		want string
	}{
		{name: "simple function example", want: "示例1 基础函数"},
		{name: "multiple returns example", want: "示例2 多返回值与命名返回"},
		{name: "variadic example", want: "示例3 可变参数"},
		{name: "closure example", want: "示例4 闭包"},
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
	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read output: %v", err)
	}

	return buf.String()
}
