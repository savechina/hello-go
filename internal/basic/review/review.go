package review

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"hello/internal/basic/review/support/tag"
	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "review", Run)
}

func main() {
	Run()
}

type notebook[T comparable] struct {
	items []T
}

func (n *notebook[T]) Add(item T) {
	if !n.Contains(item) {
		n.items = append(n.items, item)
	}
}

func (n *notebook[T]) Contains(target T) bool {
	for _, item := range n.items {
		if item == target {
			return true
		}
	}
	return false
}

type learner struct {
	Name      string
	Completed int
}

func (l *learner) Finish(topic string) string {
	if l == nil {
		return "nil learner"
	}
	l.Completed++
	return fmt.Sprintf("%s finished %s", l.Name, topic)
}

func buildStudyLog(name string, completed int) string {
	var buffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("study summary", "learner", name, "completed", completed)
	return strings.TrimSpace(buffer.String())
}

func summarizeTopics(items []string) string {
	var notes notebook[string]
	for _, item := range items {
		notes.Add(item)
	}
	return fmt.Sprintf("topics=%s count=%d", strings.Join(notes.items, ","), len(notes.items))
}

func summarizeLearnerProgress(name string, topics []string) string {
	student := &learner{Name: name}
	updates := make([]string, 0, len(topics))
	for _, topic := range topics {
		updates = append(updates, student.Finish(topic))
	}
	return fmt.Sprintf("%s %s", tag.Prefix("progress"), strings.Join(updates, " | "))
}

func summaryExample(name string, topics []string) string {
	student := &learner{Name: name}
	var notes notebook[string]
	for _, topic := range topics {
		notes.Add(topic)
		student.Finish(topic)
	}
	return fmt.Sprintf("%s %s %s", tag.Prefix("all"), summarizeTopics(notes.items), buildStudyLog(student.Name, student.Completed))
}

// Run prints the review chapter examples.
func Run() {
	fmt.Println("[review] example 1:", summarizeTopics([]string{"generics", "packages", "generics"}))
	fmt.Println("[review] example 2:", summarizeLearnerProgress("gopher", []string{"pointers", "logging"}))
	fmt.Println("[review] example 3:", summaryExample("learner", []string{"generics", "packages", "pointers", "logging"}))
}
