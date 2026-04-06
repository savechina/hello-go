package datatype

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestDatatypeExamples(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		contains []string
	}{
		{
			name:     "numeric types",
			output:   numericTypesExample(),
			contains: []string{"count=42", "price=19.95", "active=true", "label=Go 1.24"},
		},
		{
			name:     "slice operations",
			output:   sliceOperationsExample(),
			contains: []string{"scores=[80 85 90]", "window=[85 90]", "len=3"},
		},
		{
			name:     "map crud",
			output:   mapCRUDExample(),
			contains: []string{"alice=21", "insertedBob=18", "hasBob=false", "len=1"},
		},
		{
			name:     "time value",
			output:   timeValueExample(),
			contains: []string{"2026-04-05T14:30:00Z", "2026-04-07T14:30:00Z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, want := range tt.contains {
				if !strings.Contains(tt.output, want) {
					t.Fatalf("expected %q to contain %q", tt.output, want)
				}
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
		{name: "numeric example", want: "示例1 数值与基础类型"},
		{name: "slice example", want: "示例2 切片操作"},
		{name: "map example", want: "示例3 map CRUD"},
		{name: "time example", want: "示例4 时间类型"},
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
