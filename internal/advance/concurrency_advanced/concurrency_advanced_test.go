package concurrency_advanced

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestMutexCounter(t *testing.T) {
	var mu sync.Mutex
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()

	if counter != 100 {
		t.Errorf("expected 100, got %d", counter)
	}
}

func TestRWMutex(t *testing.T) {
	var rwmu sync.RWMutex
	data := map[string]int{"a": 1}

	rwmu.RLock()
	val := data["a"]
	rwmu.RUnlock()

	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}

	rwmu.Lock()
	data["b"] = 2
	rwmu.Unlock()

	rwmu.RLock()
	val2 := data["b"]
	rwmu.RUnlock()

	if val2 != 2 {
		t.Errorf("expected 2, got %d", val2)
	}
}

func TestAtomicCounter(t *testing.T) {
	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()

	if got := atomic.LoadInt64(&counter); got != 1000 {
		t.Errorf("expected 1000, got %d", got)
	}
}
