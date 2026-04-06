package visibility

import (
	"fmt"
	"strings"
)

// Profile demonstrates exported and unexported fields.
type Profile struct {
	Name  string
	score int
}

// NewProfile constructs a profile while keeping score hidden.
func NewProfile(name string, score int) Profile {
	return Profile{Name: normalizeName(name), score: score}
}

// PublicSummary exposes a safe description without exposing internals directly.
func (p Profile) PublicSummary() string {
	return fmt.Sprintf("%s(%s)", p.Name, p.ScoreBand())
}

// ScoreBand converts the hidden score into a public label.
func (p Profile) ScoreBand() string {
	switch {
	case p.score >= 90:
		return "excellent"
	case p.score >= 60:
		return "pass"
	default:
		return "retry"
	}
}

func normalizeName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "anonymous"
	}
	return strings.Title(trimmed)
}
