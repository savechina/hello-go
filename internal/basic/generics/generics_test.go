package generics

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSumValues(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{name: "three numbers", got: sumValues([]int{1, 2, 3}), want: 6},
		{name: "empty slice", got: sumValues([]int{}), want: 0},
		{name: "negative numbers", got: sumValues([]int{-2, 5, 7}), want: 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("sumValues() = %d, want %d", tt.got, tt.want)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		target string
		want   bool
	}{
		{name: "present", values: []string{"go", "package", "pointer"}, target: "package", want: true},
		{name: "missing", values: []string{"go", "package", "pointer"}, target: "generic", want: false},
		{name: "empty", values: []string{}, target: "go", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.values, tt.target); got != tt.want {
				t.Fatalf("contains() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestMapSlice(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []string
	}{
		{name: "basic mapping", values: []int{1, 2, 3}, want: []string{"n=1", "n=2", "n=3"}},
		{name: "empty mapping", values: []int{}, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapSlice(tt.values, func(value int) string { return fmt.Sprintf("n=%d", value) })
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("mapSlice() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestStackPop(t *testing.T) {
	tests := []struct {
		name    string
		values  []string
		wantTop string
		wantLen int
		wantOK  bool
	}{
		{name: "pop last value", values: []string{"a", "b", "c"}, wantTop: "c", wantLen: 2, wantOK: true},
		{name: "pop empty stack", values: []string{}, wantTop: "", wantLen: 0, wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s stack[string]
			for _, value := range tt.values {
				s.Push(value)
			}

			gotTop, gotOK := s.Pop()
			if gotTop != tt.wantTop || gotOK != tt.wantOK || s.Len() != tt.wantLen {
				t.Fatalf("Pop() = (%q, %t) len=%d, want (%q, %t) len=%d", gotTop, gotOK, s.Len(), tt.wantTop, tt.wantOK, tt.wantLen)
			}
		})
	}
}

func TestAverageValues(t *testing.T) {
	tests := []struct {
		name   string
		values []score
		want   float64
	}{
		{name: "score average", values: []score{80, 90, 100}, want: 90},
		{name: "empty average", values: []score{}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := averageValues(tt.values); got != tt.want {
				t.Fatalf("averageValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
