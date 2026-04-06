package smartpointers

import (
	"strings"
	"testing"
)

func TestReferenceCounterLifecycle(t *testing.T) {
	tests := []struct {
		name       string
		operations []string
		wantRefs   int
		wantMarked bool
	}{
		{name: "retain and release to zero", operations: []string{"add", "release", "release"}, wantRefs: 0, wantMarked: true},
		{name: "extra release stays at zero", operations: []string{"release", "release"}, wantRefs: 0, wantMarked: true},
		{name: "keeps live reference", operations: []string{"add", "release"}, wantRefs: 1, wantMarked: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counter := newRefCounter("asset")
			for _, op := range tt.operations {
				switch op {
				case "add":
					counter.AddRef()
				case "release":
					counter.Release()
				}
			}

			if counter.refs != tt.wantRefs || counter.released != tt.wantMarked {
				t.Fatalf("counter refs=%d released=%t, want refs=%d released=%t", counter.refs, counter.released, tt.wantRefs, tt.wantMarked)
			}
		})
	}
}

func TestPoolReuse(t *testing.T) {
	tests := []struct {
		name         string
		tasks        []string
		wantCreated  int
		wantContains []string
	}{
		{name: "single object reused", tasks: []string{"a", "b", "c"}, wantCreated: 1, wantContains: []string{"object#1", "task=c"}},
		{name: "empty tasks creates none", tasks: nil, wantCreated: 0, wantContains: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := simulatePoolReuse(tt.tasks)
			if !strings.Contains(got, "created=") {
				t.Fatalf("simulatePoolReuse() = %q, want created summary", got)
			}
			if tt.wantCreated == 0 && got != "created=0 " {
				t.Fatalf("simulatePoolReuse() = %q, want empty pool summary", got)
			}
			if tt.wantCreated > 0 && !strings.Contains(got, "created=1") {
				t.Fatalf("simulatePoolReuse() = %q, want created=%d", got, tt.wantCreated)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Fatalf("simulatePoolReuse() = %q, missing %q", got, want)
				}
			}
		})
	}
}

func TestProcessWithCleanup(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{name: "join values", parts: []string{"alpha", "beta"}, want: "joined=alpha/beta cleanup-items=2"},
		{name: "empty values", parts: nil, want: "joined= cleanup-items=0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processWithCleanup(tt.parts); got != tt.want {
				t.Fatalf("processWithCleanup() = %q, want %q", got, tt.want)
			}
		})
	}
}
