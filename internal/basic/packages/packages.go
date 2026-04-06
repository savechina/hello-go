package packages

import (
	"fmt"
	"strings"

	"hello/internal/basic/packages/demo/beta"
	"hello/internal/basic/packages/demo/trace"
	"hello/internal/basic/packages/demo/visibility"
	"hello/internal/chapters"
)

const moduleName = "hello"

var initOrder []string

func init() {
	chapters.Register("basic", "packages", Run)
	trace.Record("main.init")
	initOrder = trace.Events()
}

func main() {
	Run()
}

func describeImportUsage(name string, score int) string {
	profile := visibility.NewProfile(name, score)
	return fmt.Sprintf("%s | %s", beta.Description(), profile.PublicSummary())
}

func observeInitOrder() []string {
	result := make([]string, len(initOrder))
	copy(result, initOrder)
	return result
}

func explainGoModBasics() []string {
	return []string{
		"module path: " + moduleName,
		"local import path: hello/internal/basic/packages/demo/visibility",
		"go.mod defines the module root used by imports",
		"joined path: " + strings.Join([]string{moduleName, "internal", "basic", "packages", "demo", "beta"}, "/"),
	}
}

// Run prints the packages chapter examples.
func Run() {
	fmt.Println("[packages] example 1:", describeImportUsage("gopher", 88))
	fmt.Println("[packages] example 2:", strings.Join(observeInitOrder(), " -> "))
	fmt.Println("[packages] example 3:", strings.Join(explainGoModBasics(), " | "))
}
