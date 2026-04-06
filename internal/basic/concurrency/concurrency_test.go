package concurrency

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestConcurrencyHelpers(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "goroutine and channel",
			got:  fmt.Sprintf("%v", generateSequence(4)),
			want: "[1 2 3 4]",
		},
		{
			name: "worker pool order preserved",
			got:  fmt.Sprintf("%v", squareJobs([]int{2, 3, 4}, 2)),
			want: "[4 9 16]",
		},
		{
			name: "select timeout success",
			got:  waitForSignal(1*time.Millisecond, 10*time.Millisecond),
			want: "work finished",
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

func TestWaitForSignalTimeout(t *testing.T) {
	tests := []struct {
		name         string
		produceDelay time.Duration
		timeout      time.Duration
		want         string
	}{
		{
			name:         "timeout path",
			produceDelay: 10 * time.Millisecond,
			timeout:      1 * time.Millisecond,
			want:         "timeout reached",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := waitForSignal(tt.produceDelay, tt.timeout); got != tt.want {
				t.Fatalf("unexpected timeout result: got %q want %q", got, tt.want)
			}
		})
	}
}

func TestRunPrintsConcurrencyExamples(t *testing.T) {
	output := captureOutput(t, Run)

	checks := []string{
		"1) goroutine + channel => [1 2 3 4]",
		"2) worker pattern + WaitGroup => [4 9 16]",
		"3) select timeout => work finished",
	}

	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}

	if got := squareJobs([]int{1, 2, 3}, 3); !reflect.DeepEqual(got, []int{1, 4, 9}) {
		t.Fatalf("unexpected ordered worker results: %v", got)
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
