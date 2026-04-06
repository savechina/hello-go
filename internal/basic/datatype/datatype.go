package datatype

import (
	"fmt"
	"time"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "datatype", Run)
}

// Run prints runnable examples for core Go data types.
func Run() {
	for _, line := range renderLines() {
		fmt.Println(line)
	}
}

func renderLines() []string {
	return []string{
		numericTypesExample(),
		sliceOperationsExample(),
		mapCRUDExample(),
		timeValueExample(),
	}
}

func numericTypesExample() string {
	count := 42
	price := 19.95
	active := true
	label := "Go 1.24"

	return fmt.Sprintf(
		"示例1 数值与基础类型: count=%d price=%.2f active=%t label=%s",
		count,
		price,
		active,
		label,
	)
}

func sliceOperationsExample() string {
	scores := []int{80, 85}
	scores = append(scores, 90)
	window := scores[1:]

	return fmt.Sprintf(
		"示例2 切片操作: scores=%v window=%v len=%d cap=%d",
		scores,
		window,
		len(scores),
		cap(scores),
	)
}

func mapCRUDExample() string {
	ages := map[string]int{"Alice": 20}
	ages["Bob"] = 18
	insertedBob := ages["Bob"]
	ages["Alice"] = 21
	delete(ages, "Bob")
	_, hasBob := ages["Bob"]

	return fmt.Sprintf(
		"示例3 map CRUD: alice=%d insertedBob=%d hasBob=%t len=%d",
		ages["Alice"],
		insertedBob,
		hasBob,
		len(ages),
	)
}

func timeValueExample() string {
	createdAt := time.Date(2026, time.April, 5, 14, 30, 0, 0, time.UTC)
	deadline := createdAt.Add(48 * time.Hour)

	return fmt.Sprintf(
		"示例4 时间类型: createdAt=%s deadline=%s",
		createdAt.Format(time.RFC3339),
		deadline.Format(time.RFC3339),
	)
}
