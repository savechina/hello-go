package concurrency_advanced

import (
	"fmt"
	"sync"
	"sync/atomic"

	"hello/internal/chapters"
)

// Run demonstrates advanced concurrency: Mutex, RWMutex, and atomic operations.
func Run() {
	fmt.Println("=== 高级并发 (Advanced Concurrency) ===")

	// 示例1: sync.Mutex 互斥锁 (Mutual exclusion)
	var mu sync.Mutex
	counter := 0
	var wg1 sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg1.Add(1)
		go func() {
			defer wg1.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg1.Wait()
	fmt.Printf("  示例1: Mutex 计数 = %d (预期 100)\n", counter)

	// 示例2: sync.RWMutex 读写锁 (Read-write lock)
	var rwmu sync.RWMutex
	data := map[string]int{"a": 1}
	var wg2 sync.WaitGroup
	// 写操作
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		rwmu.Lock()
		data["b"] = 2
		rwmu.Unlock()
	}()
	// 读操作 (可并发)
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			rwmu.RLock()
			_ = data["a"]
			rwmu.RUnlock()
		}()
	}
	wg2.Wait()
	fmt.Printf("  示例2: RWMutex 数据 = %v\n", data)

	// 示例3: atomic 原子操作 (Atomic operations)
	var atomicCounter int64
	var wg3 sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg3.Add(1)
		go func() {
			defer wg3.Done()
			atomic.AddInt64(&atomicCounter, 1)
		}()
	}
	wg3.Wait()
	fmt.Printf("  示例3: atomic 计数 = %d (预期 1000)\n", atomic.LoadInt64(&atomicCounter))
}

func init() {
	chapters.Register("advance", "concurrency_advanced", Run)
}
