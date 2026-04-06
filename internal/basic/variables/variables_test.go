package variables

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestVariableExamples(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		contains []string
	}{
		{
			name:     "basic var",
			output:   basicVarExample(),
			contains: []string{"language=Go", "lessonCount=12", "ready=true"},
		},
		{
			name:     "const iota",
			output:   constIotaExample(),
			contains: []string{"course=hello-go", "draft=0", "review=1", "published=2"},
		},
		{
			name:     "short declaration",
			output:   shortDeclarationExample(),
			contains: []string{"name=Alice", "chapter=variables", "completed=true"},
		},
		{
			name:     "type inference",
			output:   typeInferenceExample(),
			contains: []string{"total=3(int)", "progress=75.5(float64)", "note=type inference(string)"},
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
		{name: "basic example", want: "示例1 基础变量"},
		{name: "const example", want: "示例2 常量与iota"},
		{name: "short declaration example", want: "示例3 短变量声明"},
		{name: "type inference example", want: "示例4 类型推断"},
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
