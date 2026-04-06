package alpha

import "hello/internal/basic/packages/demo/trace"

func init() {
	trace.Record("alpha.init")
}

// Title returns the alpha package label.
func Title() string {
	return "alpha"
}
