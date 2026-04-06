package interfaces

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestInterfaceHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "implicit implementation",
			got:  announceGreeter(robot{name: "R2"}),
			want: "R2 says hello",
		},
		{
			name: "writer pattern",
			got:  captureLogLine("info", "interfaces keep behavior separate"),
			want: "INFO: interfaces keep behavior separate",
		},
		{
			name: "type switch",
			got:  inspectValues([]any{"go", 7, robot{name: "Mika"}}),
			want: "string => GO | int => 14 | greeter => Mika says hello",
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

func TestRunPrintsInterfaceExamples(t *testing.T) {
	output := captureOutput(t, Run)

	tests := []string{
		"1) implicit implementation => R2 says hello",
		"2) io.Writer pattern => INFO: interfaces keep behavior separate",
		"3) empty interface + type switch => string => GO | int => 14 | greeter => Mika says hello",
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
