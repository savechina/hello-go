package errorhandling

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestErrorHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "basic sentinel error",
			got:  validateAmount(-1).Error(),
			want: "amount must be positive",
		},
		{
			name: "wrapped error summary",
			got: func() string {
				_, err := lookupSetting(map[string]string{"mode": "dev"}, "timeout")
				return summarizeError(err)
			}(),
			want: "missing setting detected",
		},
		{
			name: "errors.As custom error",
			got: func() string {
				_, err := parseRetryCount("abc")
				return summarizeError(err)
			}(),
			want: "field error on retry count with value \"abc\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("unexpected helper output: got %q want %q", tt.got, tt.want)
			}
		})
	}
}

func TestRunPrintsErrorExamples(t *testing.T) {
	output := captureOutput(t, Run)

	tests := []string{
		"1) errors.New => amount must be positive",
		"2) fmt.Errorf + %w => lookup \"timeout\": setting not found",
		"3) errors.Is / errors.As => missing setting detected | field error on retry count with value \"abc\"",
	}

	for _, want := range tests {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, r); err != nil {
		t.Fatalf("copy output: %v", err)
	}

	return buffer.String()
}
