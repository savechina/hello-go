package logging

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "logging", Run)
}

func main() {
	Run()
}

type handlerState struct {
	records []string
}

type memoryHandler struct {
	level slog.Leveler
	attrs []slog.Attr
	group string
	state *handlerState
}

func newMemoryHandler(level slog.Leveler) *memoryHandler {
	return &memoryHandler{level: level, state: &handlerState{}}
}

func (h *memoryHandler) Enabled(_ context.Context, level slog.Level) bool {
	if h.level == nil {
		return true
	}
	return level >= h.level.Level()
}

func (h *memoryHandler) Handle(_ context.Context, record slog.Record) error {
	parts := []string{
		"level=" + record.Level.String(),
		"msg=" + record.Message,
	}

	for _, attr := range h.attrs {
		parts = append(parts, formatAttr(h.group, attr))
	}

	record.Attrs(func(attr slog.Attr) bool {
		parts = append(parts, formatAttr(h.group, attr))
		return true
	})

	h.state.records = append(h.state.records, strings.Join(parts, " "))
	return nil

}

func (h *memoryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	clone := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	clone = append(clone, h.attrs...)
	clone = append(clone, attrs...)
	return &memoryHandler{level: h.level, attrs: clone, group: h.group, state: h.state}
}

func (h *memoryHandler) WithGroup(name string) slog.Handler {
	nextGroup := name
	if h.group != "" {
		nextGroup = h.group + "." + name
	}
	return &memoryHandler{level: h.level, attrs: append([]slog.Attr{}, h.attrs...), group: nextGroup, state: h.state}
}

func (h *memoryHandler) Records() []string {
	result := make([]string, len(h.state.records))
	copy(result, h.state.records)
	return result
}

func formatAttr(group string, attr slog.Attr) string {
	key := attr.Key
	if group != "" {
		key = group + "." + key
	}
	return fmt.Sprintf("%s=%v", key, attr.Value.Any())
}

func basicLogOutput(topic string) string {
	var buffer bytes.Buffer
	logger := log.New(&buffer, "basic ", 0)
	logger.Println("studying", topic)
	return strings.TrimSpace(buffer.String())
}

func structuredLogOutput(orderID string, amount float64) string {
	var buffer bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("order created", "order_id", orderID, "amount", amount)
	return strings.TrimSpace(buffer.String())
}

func customHandlerOutput(minLevel slog.Level, module string) []string {
	levelVar := new(slog.LevelVar)
	levelVar.Set(minLevel)
	handler := newMemoryHandler(levelVar)
	logger := slog.New(handler).With("module", module)
	logger.Info("skip info")
	logger.Warn("keep warn", "attempt", 2)
	logger.Error("keep error", "attempt", 3)
	return handler.Records()
}

// Run prints the logging chapter examples.
func Run() {
	fmt.Println("[logging] example 1:", basicLogOutput("log package"))
	fmt.Println("[logging] example 2:", structuredLogOutput("A-100", 19.9))
	fmt.Println("[logging] example 3:", strings.Join(customHandlerOutput(slog.LevelWarn, "study"), " | "))
}
