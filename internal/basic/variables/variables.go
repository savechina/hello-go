package variables

import (
	"fmt"

	"hello/internal/chapters"
)

const (
	stageDraft = iota
	stageReview
	stagePublished
)

func init() {
	chapters.Register("basic", "variables", Run)
}

// Run prints runnable examples for variables, constants, and type inference.
func Run() {
	for _, line := range renderLines() {
		fmt.Println(line)
	}
}

func renderLines() []string {
	return []string{
		basicVarExample(),
		constIotaExample(),
		shortDeclarationExample(),
		typeInferenceExample(),
	}
}

func basicVarExample() string {
	var language string = "Go"
	var lessonCount int = 12
	var ready bool = true

	return fmt.Sprintf(
		"示例1 基础变量: language=%s lessonCount=%d ready=%t",
		language,
		lessonCount,
		ready,
	)
}

func constIotaExample() string {
	const courseName = "hello-go"

	return fmt.Sprintf(
		"示例2 常量与iota: course=%s draft=%d review=%d published=%d",
		courseName,
		stageDraft,
		stageReview,
		stagePublished,
	)
}

func shortDeclarationExample() string {
	name, chapter, completed := "Alice", "variables", false
	completed = true

	return fmt.Sprintf(
		"示例3 短变量声明: name=%s chapter=%s completed=%t",
		name,
		chapter,
		completed,
	)
}

func typeInferenceExample() string {
	total := 3
	progress := 75.5
	note := "type inference"

	return fmt.Sprintf(
		"示例4 类型推断: total=%d(%T) progress=%.1f(%T) note=%s(%T)",
		total,
		total,
		progress,
		progress,
		note,
		note,
	)
}
