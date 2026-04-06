package interfaces

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "interfaces", Run)
}

func main() {
	Run()
}

// Run executes the interfaces chapter examples.
func Run() {
	examples := []string{
		"1) implicit implementation => " + announceGreeter(robot{name: "R2"}),
		"2) io.Writer pattern => " + captureLogLine("INFO", "interfaces keep behavior separate"),
		"3) empty interface + type switch => " + inspectValues([]any{"go", 7, robot{name: "Mika"}}),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

type greeter interface {
	greet() string
}

type robot struct {
	name string
}

func (r robot) greet() string {
	return fmt.Sprintf("%s says hello", r.name)
}

func announceGreeter(g greeter) string {
	return g.greet()
}

func writeLogLine(writer io.Writer, level string, message string) error {
	_, err := fmt.Fprintf(writer, "%s: %s", strings.ToUpper(level), message)
	return err
}

func captureLogLine(level string, message string) string {
	var buffer bytes.Buffer
	if err := writeLogLine(&buffer, level, message); err != nil {
		return err.Error()
	}

	return buffer.String()
}

func inspectValue(value any) string {
	switch typed := value.(type) {
	case string:
		return fmt.Sprintf("string => %s", strings.ToUpper(typed))
	case int:
		number, ok := value.(int)
		if !ok {
			return "int assertion failed"
		}
		return fmt.Sprintf("int => %d", number*2)
	case greeter:
		return fmt.Sprintf("greeter => %s", typed.greet())
	default:
		return fmt.Sprintf("unknown => %T", value)
	}
}

func inspectValues(values []any) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, inspectValue(value))
	}

	return strings.Join(parts, " | ")
}
