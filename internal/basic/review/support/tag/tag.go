package tag

import "strings"

// Prefix builds a review label for output.
func Prefix(section string) string {
	return "[review:" + strings.ToUpper(strings.TrimSpace(section)) + "]"
}
