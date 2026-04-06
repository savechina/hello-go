package structs

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestStructHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "struct literal summary",
			got:  buildProfile("Alice", 30, "Taipei"),
			want: "Alice is 30 years old and lives in Taipei",
		},
		{
			name: "method updates age",
			got:  celebrateBirthday("Bob", 27),
			want: "Bob is 28 years old and lives in Taichung",
		},
		{
			name: "embedding reuses profile methods",
			got:  describePromotion("Carol", 32, "Kaohsiung", "Platform", "Engineer"),
			want: "Carol [Platform] Senior Engineer -> Carol is 32 years old and lives in Kaohsiung",
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

func TestRunPrintsStructExamples(t *testing.T) {
	output := captureOutput(t, Run)

	tests := []string{
		"1) struct literals => Alice is 30 years old and lives in Taipei",
		"2) methods => Bob is 28 years old and lives in Taichung",
		"3) embedding => Carol [Platform] Senior Engineer -> Carol is 32 years old and lives in Kaohsiung",
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
