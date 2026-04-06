package logging

import (
	"log/slog"
	"strings"
	"testing"
)

func TestBasicLogOutput(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		want  []string
	}{
		{name: "basic logger", topic: "package", want: []string{"basic", "studying package"}},
		{name: "different topic", topic: "slog", want: []string{"studying slog"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := basicLogOutput(tt.topic)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("basicLogOutput() = %q, want substring %q", got, want)
				}
			}
		})
	}
}

func TestStructuredLogOutput(t *testing.T) {
	tests := []struct {
		name    string
		orderID string
		amount  float64
		want    []string
	}{
		{name: "text handler output", orderID: "A-100", amount: 19.9, want: []string{"level=INFO", "msg=\"order created\"", "order_id=A-100", "amount=19.9"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := structuredLogOutput(tt.orderID, tt.amount)
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("structuredLogOutput() = %q, want substring %q", got, want)
				}
			}
		})
	}
}

func TestCustomHandlerOutput(t *testing.T) {
	tests := []struct {
		name     string
		minLevel slog.Level
		wantLen  int
		wants    []string
	}{
		{name: "warn and error kept", minLevel: slog.LevelWarn, wantLen: 2, wants: []string{"level=WARN", "level=ERROR", "module=study"}},
		{name: "info and above kept", minLevel: slog.LevelInfo, wantLen: 3, wants: []string{"level=INFO", "level=WARN", "level=ERROR"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := customHandlerOutput(tt.minLevel, "study")
			if len(got) != tt.wantLen {
				t.Fatalf("customHandlerOutput() len = %d, want %d", len(got), tt.wantLen)
			}

			joined := strings.Join(got, " | ")
			for _, want := range tt.wants {
				if !strings.Contains(joined, want) {
					t.Fatalf("customHandlerOutput() = %q, want substring %q", joined, want)
				}
			}
		})
	}
}
