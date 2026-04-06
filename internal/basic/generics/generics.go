package generics

import (
	"fmt"
	"strings"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "generics", Run)
}

func main() {
	Run()
}

type number interface {
	~int | ~int64 | ~float64
}

type score float64

type stack[T any] struct {
	items []T
}

func (s *stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	index := len(s.items) - 1
	item := s.items[index]
	s.items = s.items[:index]
	return item, true
}

func (s *stack[T]) Len() int {
	return len(s.items)
}

func sumValues[T number](values []T) T {
	var total T
	for _, value := range values {
		total += value
	}
	return total
}

func averageValues[T number](values []T) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(sumValues(values)) / float64(len(values))
}

func mapSlice[T any, R any](values []T, mapper func(T) R) []R {
	result := make([]R, 0, len(values))
	for _, value := range values {
		result = append(result, mapper(value))
	}
	return result
}

func contains[T comparable](values []T, target T) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func exampleGenericFunction() string {
	numbers := []int{2, 4, 6}
	total := sumValues(numbers)
	doubled := mapSlice(numbers, func(value int) string {
		return fmt.Sprintf("%d->%d", value, value*2)
	})
	return fmt.Sprintf("sum=%d mapped=%s", total, strings.Join(doubled, ", "))
}

func exampleComparableConstraint() string {
	keywords := []string{"go", "generic", "comparable"}
	return fmt.Sprintf("contains generic=%t contains interface=%t", contains(keywords, "generic"), contains(keywords, "interface"))
}

func exampleCustomConstraintAndType() string {
	grades := []score{90.5, 88.0, 95.5}
	var lessons stack[string]
	lessons.Push("type parameters")
	lessons.Push("constraints")
	top, _ := lessons.Pop()
	return fmt.Sprintf("average=%.1f next=%s remaining=%d", averageValues(grades), top, lessons.Len())
}

// Run prints the generics chapter examples.
func Run() {
	fmt.Println("[generics] example 1:", exampleGenericFunction())
	fmt.Println("[generics] example 2:", exampleComparableConstraint())
	fmt.Println("[generics] example 3:", exampleCustomConstraintAndType())
}
