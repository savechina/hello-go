package functions

import (
	"fmt"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "functions", Run)
}

// Run prints runnable examples for function definitions and common patterns.
func Run() {
	for _, line := range renderLines() {
		fmt.Println(line)
	}
}

func renderLines() []string {
	return []string{
		simpleFunctionExample(),
		multipleAndNamedReturnExample(),
		variadicExample(),
		closureExample(),
	}
}

func simpleFunctionExample() string {
	return fmt.Sprintf("示例1 基础函数: %s", greet("Gopher"))
}

func multipleAndNamedReturnExample() string {
	quotient, remainder, err := divideAndRemainder(17, 5)
	if err != nil {
		return fmt.Sprintf("示例2 多返回值与命名返回: error=%v", err)
	}

	area, perimeter := rectangleMetrics(3, 4)

	return fmt.Sprintf(
		"示例2 多返回值与命名返回: quotient=%d remainder=%d area=%.1f perimeter=%.1f",
		quotient,
		remainder,
		area,
		perimeter,
	)
}

func variadicExample() string {
	return fmt.Sprintf("示例3 可变参数: total=%d", sumAll(3, 5, 7))
}

func closureExample() string {
	counter := makeCounter(0)

	return fmt.Sprintf(
		"示例4 闭包: first=%d second=%d third=%d",
		counter(),
		counter(),
		counter(),
	)
}

func greet(name string) string {
	return "Hello, " + name
}

func divideAndRemainder(total int, group int) (quotient int, remainder int, err error) {
	if group == 0 {
		return 0, 0, fmt.Errorf("group must not be zero")
	}

	quotient = total / group
	remainder = total % group

	return quotient, remainder, nil
}

func sumAll(nums ...int) int {
	total := 0
	for _, num := range nums {
		total += num
	}

	return total
}

func rectangleMetrics(width float64, height float64) (area float64, perimeter float64) {
	area = width * height
	perimeter = 2 * (width + height)

	return area, perimeter
}

func makeCounter(start int) func() int {
	current := start

	return func() int {
		current++
		return current
	}
}
