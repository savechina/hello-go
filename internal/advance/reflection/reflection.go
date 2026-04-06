package reflection

import (
	"fmt"
	"reflect"
	"strings"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("advance", "reflection", Run)
}

type taggedUser struct {
	Name  string `json:"name" db:"user_name"`
	Level int    `json:"level" db:"user_level"`
}

type greeter struct {
	Prefix string
}

func (g greeter) Greet(name string) string {
	return fmt.Sprintf("%s, %s", g.Prefix, name)
}

func describeValue(input any) string {
	value := reflect.ValueOf(input)
	typ := reflect.TypeOf(input)
	if !value.IsValid() || typ == nil {
		return "invalid value"
	}

	return fmt.Sprintf("type=%s kind=%s value=%v", typ.String(), value.Kind(), value.Interface())
}

func readStructTags(input any) string {
	typ := reflect.TypeOf(input)
	if typ == nil {
		return "no tags"
	}
	if typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return "no tags"
	}

	parts := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		parts = append(parts, fmt.Sprintf("%s json=%s db=%s", field.Name, field.Tag.Get("json"), field.Tag.Get("db")))
	}

	return strings.Join(parts, " | ")
}

func callMethod(target any, method string, args ...string) string {
	value := reflect.ValueOf(target)
	selected := value.MethodByName(method)
	if !selected.IsValid() {
		return "method not found"
	}

	inputs := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		inputs = append(inputs, reflect.ValueOf(arg))
	}

	outputs := selected.Call(inputs)
	if len(outputs) == 0 {
		return "no result"
	}

	return fmt.Sprint(outputs[0].Interface())
}

// Run prints the reflection examples.
func Run() {
	fmt.Println("[reflection] example 1:", describeValue(taggedUser{Name: "gopher", Level: 3}))
	fmt.Println("[reflection] example 2:", readStructTags(taggedUser{}))
	fmt.Println("[reflection] example 3:", callMethod(greeter{Prefix: "hello"}, "Greet", "Go"))
}
