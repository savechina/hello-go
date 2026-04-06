package context

import (
	"context"
	"testing"
	"time"
)

func TestContextWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	select {
	case <-ctx.Done():
		if ctx.Err() != context.Canceled {
			t.Errorf("expected Canceled, got %v", ctx.Err())
		}
	default:
		t.Error("context should be done after cancel")
	}
}

func TestContextWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	time.Sleep(20 * time.Millisecond)

	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("expected DeadlineExceeded, got %v", ctx.Err())
		}
	default:
		t.Error("context should be done after timeout")
	}
}

func TestContextWithDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Millisecond))
	defer cancel()

	time.Sleep(20 * time.Millisecond)

	select {
	case <-ctx.Done():
		if ctx.Err() != context.DeadlineExceeded {
			t.Errorf("expected DeadlineExceeded, got %v", ctx.Err())
		}
	default:
		t.Error("context should be done after deadline")
	}
}

func TestContextValuePropagation(t *testing.T) {
	type key string
	const testKey key = "test"

	ctx := context.WithValue(context.Background(), testKey, "hello")

	if got := ctx.Value(testKey); got != "hello" {
		t.Errorf("expected 'hello', got %v", got)
	}
}
