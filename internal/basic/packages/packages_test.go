package packages

import (
	"reflect"
	"strings"
	"testing"
)

func TestDescribeImportUsage(t *testing.T) {
	tests := []struct {
		name  string
		input string
		score int
		wants []string
	}{
		{name: "normal name", input: "gopher", score: 88, wants: []string{"beta imports alpha", "Gopher(pass)"}},
		{name: "blank name", input: " ", score: 50, wants: []string{"anonymous", "retry"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := describeImportUsage(tt.input, tt.score)
			for _, want := range tt.wants {
				if !strings.Contains(strings.ToLower(got), strings.ToLower(want)) {
					t.Fatalf("describeImportUsage() = %q, want substring %q", got, want)
				}
			}
		})
	}
}

func TestObserveInitOrder(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{name: "stable import order", want: []string{"alpha.init", "beta.init", "main.init"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := observeInitOrder(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("observeInitOrder() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExplainGoModBasics(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{name: "module summary", want: []string{"module path: hello", "local import path", "joined path: hello/internal/basic/packages/demo/beta"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := strings.Join(explainGoModBasics(), " | ")
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("explainGoModBasics() = %q, want substring %q", got, want)
				}
			}
		})
	}
}
