package review

import (
	"strings"
	"testing"
)

func TestNotebookAddAndContains(t *testing.T) {
	tests := []struct {
		name   string
		items  []string
		target string
		want   bool
		count  int
	}{
		{name: "deduplicate items", items: []string{"go", "go", "slog"}, target: "go", want: true, count: 2},
		{name: "missing target", items: []string{"go", "pointer"}, target: "generic", want: false, count: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var notes notebook[string]
			for _, item := range tt.items {
				notes.Add(item)
			}

			if got := notes.Contains(tt.target); got != tt.want || len(notes.items) != tt.count {
				t.Fatalf("Contains()=%t len=%d, want %t len=%d", got, len(notes.items), tt.want, tt.count)
			}
		})
	}
}

func TestLearnerFinish(t *testing.T) {
	tests := []struct {
		name      string
		learner   *learner
		topic     string
		wantText  string
		wantCount int
	}{
		{name: "real learner", learner: &learner{Name: "gopher"}, topic: "generics", wantText: "gopher finished generics", wantCount: 1},
		{name: "nil learner", learner: nil, topic: "generics", wantText: "nil learner", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.learner.Finish(tt.topic)
			if got != tt.wantText {
				t.Fatalf("Finish() = %q, want %q", got, tt.wantText)
			}
			if tt.learner != nil && tt.learner.Completed != tt.wantCount {
				t.Fatalf("Completed = %d, want %d", tt.learner.Completed, tt.wantCount)
			}
		})
	}
}

func TestSummaryExample(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		wants []string
	}{
		{name: "combined summary", input: []string{"generics", "packages", "logging"}, wants: []string{"[review:ALL]", "topics=generics,packages,logging", "learner=student"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := summaryExample("student", tt.input)
			for _, want := range tt.wants {
				if !strings.Contains(got, want) {
					t.Fatalf("summaryExample() = %q, want substring %q", got, want)
				}
			}
		})
	}
}
