package testing

import (
	"bytes"
	"io"
	"os"
	"strings"
	stdtesting "testing"
)

func TestTestingHelpers(t *stdtesting.T) {
	tests := []struct {
		name string
		run  func(t *stdtesting.T)
	}{
		{
			name: "grade label",
			run: func(t *stdtesting.T) {
				cases := []struct {
					score int
					want  string
				}{
					{score: 95, want: "excellent"},
					{score: 70, want: "pass"},
					{score: 42, want: "retry"},
				}

				for _, item := range cases {
					if got := gradeLabel(item.score); got != item.want {
						t.Fatalf("score %d: want %q, got %q", item.score, item.want, got)
					}
				}
			},
		},
		{
			name: "normalize slug",
			run: func(t *stdtesting.T) {
				cases := []struct {
					input string
					want  string
				}{
					{input: " Go Testing ", want: "go-testing"},
					{input: "Fuzz__Case###", want: "fuzz-case"},
					{input: "123 Ready", want: "123-ready"},
				}

				for _, item := range cases {
					if got := normalizeSlug(item.input); got != item.want {
						t.Fatalf("input %q: want %q, got %q", item.input, item.want, got)
					}
				}
			},
		},
		{
			name: "join words helpers agree",
			run: func(t *stdtesting.T) {
				cases := []struct {
					words []string
					want  string
				}{
					{words: []string{"go", "bench"}, want: "go,bench"},
					{words: []string{"a", "b", "c"}, want: "a,b,c"},
				}

				for _, item := range cases {
					if got := joinWordsPlus(item.words); got != item.want {
						t.Fatalf("plus want %q, got %q", item.want, got)
					}
					if got := joinWordsBuilder(item.words); got != item.want {
						t.Fatalf("builder want %q, got %q", item.want, got)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}

func TestRunOutput(t *stdtesting.T) {
	output := captureOutput(t, Run)

	tests := []struct {
		name string
		want string
	}{
		{name: "table driven example", want: "示例1 表驱动测试"},
		{name: "benchmark example", want: "示例2 基准测试思维"},
		{name: "fuzz example", want: "示例3 模糊测试基础"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *stdtesting.T) {
			if !strings.Contains(output, tt.want) {
				t.Fatalf("expected output to contain %q, got %q", tt.want, output)
			}
		})
	}
}

func BenchmarkJoinWordsPlus(b *stdtesting.B) {
	words := []string{"go", "benchmark", "plus", "concatenation", "comparison"}

	for range b.N {
		_ = joinWordsPlus(words)
	}
}

func BenchmarkJoinWordsBuilder(b *stdtesting.B) {
	words := []string{"go", "benchmark", "builder", "concatenation", "comparison"}

	for range b.N {
		_ = joinWordsBuilder(words)
	}
}

func FuzzNormalizeSlug(f *stdtesting.F) {
	seeds := []string{
		"Go Testing",
		"Fuzz__Case###",
		"  Bench   Marker  ",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *stdtesting.T, input string) {
		normalized := normalizeSlug(input)
		if normalized != strings.ToLower(normalized) {
			t.Fatalf("expected lowercase slug, got %q", normalized)
		}
		if strings.Contains(normalized, " ") {
			t.Fatalf("expected no spaces, got %q", normalized)
		}
		if normalizeSlug(normalized) != normalized {
			t.Fatalf("expected idempotent normalization, got %q", normalized)
		}
	})
}

func captureOutput(t *stdtesting.T, runner func()) string {
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
