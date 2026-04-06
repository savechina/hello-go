package testing

import (
	"fmt"
	"strings"
	"unicode"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("advance", "testing", Run)
}

// Run prints runnable examples for table-driven tests, benchmarks, and fuzz testing basics.
func Run() {
	examples := []string{
		tableDrivenTestingExample(),
		benchmarkComparisonExample(),
		fuzzingBasicsExample(),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

func tableDrivenTestingExample() string {
	cases := []struct {
		score int
		want  string
	}{
		{score: 95, want: "excellent"},
		{score: 75, want: "pass"},
		{score: 50, want: "retry"},
	}

	passed := 0
	for _, item := range cases {
		if gradeLabel(item.score) == item.want {
			passed++
		}
	}

	return fmt.Sprintf("示例1 表驱动测试: cases=%d passed=%d", len(cases), passed)
}

func benchmarkComparisonExample() string {
	words := []string{"go", "benchmark", "builder", "comparison"}
	plusResult := joinWordsPlus(words)
	builderResult := joinWordsBuilder(words)

	return fmt.Sprintf(
		"示例2 基准测试思维: plus=%s builder=%s equal=%t",
		plusResult,
		builderResult,
		plusResult == builderResult,
	)
}

func fuzzingBasicsExample() string {
	seed := "  Go_Fuzzing  Basics!!!  "
	normalized := normalizeSlug(seed)

	return fmt.Sprintf("示例3 模糊测试基础: seed=%q normalized=%q", seed, normalized)
}

func gradeLabel(score int) string {
	switch {
	case score >= 90:
		return "excellent"
	case score >= 60:
		return "pass"
	default:
		return "retry"
	}
}

func joinWordsPlus(words []string) string {
	result := ""
	for index, word := range words {
		if index > 0 {
			result += ","
		}
		result += word
	}

	return result
}

func joinWordsBuilder(words []string) string {
	var builder strings.Builder
	for index, word := range words {
		if index > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(word)
	}

	return builder.String()
}

func normalizeSlug(input string) string {
	var builder strings.Builder
	lastWasDash := false

	for _, r := range strings.ToLower(strings.TrimSpace(input)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(r)
			lastWasDash = false
		case r == '-' || r == '_' || unicode.IsSpace(r):
			if builder.Len() > 0 && !lastWasDash {
				builder.WriteRune('-')
				lastWasDash = true
			}
		default:
			if builder.Len() > 0 && !lastWasDash {
				builder.WriteRune('-')
				lastWasDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-")
}
