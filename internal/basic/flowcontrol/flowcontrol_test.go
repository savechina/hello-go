package flowcontrol

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestFlowControlHelpers(t *testing.T) {
	tests := []struct {
		name      string
		checkFunc func(t *testing.T)
	}{
		{
			name: "classify score",
			checkFunc: func(t *testing.T) {
				cases := []struct {
					score int
					want  string
				}{
					{score: 95, want: "优秀"},
					{score: 75, want: "及格"},
					{score: 50, want: "继续练习"},
				}

				for _, tc := range cases {
					t.Run(tc.want, func(t *testing.T) {
						if got := classifyScore(tc.score); got != tc.want {
							t.Fatalf("expected %q, got %q", tc.want, got)
						}
					})
				}
			},
		},
		{
			name: "label day",
			checkFunc: func(t *testing.T) {
				cases := []struct {
					day  string
					want string
				}{
					{day: "Monday", want: "新的开始"},
					{day: "Sunday", want: "周末"},
					{day: "Wednesday", want: "工作日"},
				}

				for _, tc := range cases {
					t.Run(tc.day, func(t *testing.T) {
						if got := labelDay(tc.day); got != tc.want {
							t.Fatalf("expected %q, got %q", tc.want, got)
						}
					})
				}
			},
		},
		{
			name: "defer stack order",
			checkFunc: func(t *testing.T) {
				want := []string{"enter", "leave", "defer:second", "defer:first"}
				if got := collectDeferStack(); !reflect.DeepEqual(got, want) {
					t.Fatalf("expected %v, got %v", want, got)
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
		{name: "conditional example", want: "示例1 条件判断"},
		{name: "loop example", want: "示例2 for 循环"},
		{name: "switch example", want: "示例3 switch 分支"},
		{name: "defer example", want: "示例4 defer 栈"},
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
