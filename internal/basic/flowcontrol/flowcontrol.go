package flowcontrol

import (
	"fmt"
	"strings"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "flowcontrol", Run)
}

// Run prints runnable examples for conditional branches, loops, switches, and defer.
func Run() {
	for _, line := range renderLines() {
		fmt.Println(line)
	}
}

func renderLines() []string {
	return []string{
		conditionalExample(),
		loopExample(),
		switchExample(),
		deferExample(),
	}
}

func conditionalExample() string {
	return fmt.Sprintf("示例1 条件判断: score=88 result=%s", classifyScore(88))
}

func loopExample() string {
	classicTotal := 0
	for i := 1; i <= 4; i++ {
		classicTotal += i
	}

	words := []string{"go", "is", "fun"}
	rangeChars := 0
	for _, word := range words {
		rangeChars += len(word)
	}

	return fmt.Sprintf("示例2 for 循环: classicTotal=%d rangeChars=%d", classicTotal, rangeChars)
}

func switchExample() string {
	return fmt.Sprintf("示例3 switch 分支: monday=%s sunday=%s", labelDay("Monday"), labelDay("Sunday"))
}

func deferExample() string {
	return fmt.Sprintf("示例4 defer 栈: %s", strings.Join(collectDeferStack(), " -> "))
}

func classifyScore(score int) string {
	if score >= 90 {
		return "优秀"
	}

	if score >= 60 {
		return "及格"
	}

	return "继续练习"
}

func labelDay(day string) string {
	switch day {
	case "Saturday", "Sunday":
		return "周末"
	case "Monday":
		return "新的开始"
	default:
		return "工作日"
	}
}

func collectDeferStack() (steps []string) {
	steps = append(steps, "enter")
	defer func() {
		steps = append(steps, "defer:first")
	}()
	defer func() {
		steps = append(steps, "defer:second")
	}()
	steps = append(steps, "leave")

	return steps
}
