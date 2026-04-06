package context

import (
	"context"
	"fmt"
	"time"

	"hello/internal/chapters"
)

// Run demonstrates Go context usage: cancellation, timeout, and deadline.
func Run() {
	fmt.Println("=== Context 上下文 (Context) ===")

	// 示例1: 基础取消 (Basic cancellation)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-time.After(100 * time.Millisecond):
			fmt.Println("  示例1: 任务完成")
			cancel()
		case <-ctx.Done():
			fmt.Println("  示例1: 任务被取消")
		}
	}()
	cancel() // 立即取消
	time.Sleep(50 * time.Millisecond)

	// 示例2: 超时控制 (Timeout control)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  示例2: 任务完成")
	case <-ctx2.Done():
		fmt.Printf("  示例2: 超时 - %v\n", ctx2.Err())
	}

	// 示例3: 截止时间 (Deadline)
	ctx3, cancel3 := context.WithDeadline(context.Background(), time.Now().Add(200*time.Millisecond))
	defer cancel3()
	select {
	case <-time.After(300 * time.Millisecond):
		fmt.Println("  示例3: 任务完成")
	case <-ctx3.Done():
		fmt.Printf("  示例3: 截止时间到达 - %v\n", ctx3.Err())
	}
}

func init() {
	chapters.Register("advance", "context", Run)
}
