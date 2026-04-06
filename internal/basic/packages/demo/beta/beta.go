package beta

import (
	"fmt"

	"hello/internal/basic/packages/demo/alpha"
	"hello/internal/basic/packages/demo/trace"
)

func init() {
	trace.Record("beta.init")
}

// Description returns a sentence showing nested imports.
func Description() string {
	return fmt.Sprintf("beta imports %s", alpha.Title())
}
