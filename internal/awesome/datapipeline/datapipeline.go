package datapipeline

import (
	"fmt"
	"sync"
	"time"

	"hello/internal/chapters"
)

// Run demonstrates a data processing pipeline with worker pools and graceful shutdown.
func Run() {
	fmt.Println("=== 实战项目：数据处理管道 (Data Pipeline) ===")

	// 示例1: Worker Pool 工作池 (Worker pool pattern)
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("  示例1: Worker Pool 处理结果:")
	for r := range results {
		fmt.Printf("    结果: %d\n", r)
	}

	// 示例2: Graceful Shutdown 优雅关闭 (Graceful shutdown with context)
	done := make(chan struct{})
	go func() {
		fmt.Println("  示例2: 模拟长时间运行的任务...")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("  示例2: 任务完成，发送关闭信号")
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("  示例2: 收到完成信号")
	case <-time.After(200 * time.Millisecond):
		fmt.Println("  示例2: 超时，强制退出")
	}

	// 示例3: Fan-out/Fan-in 扇出扇入 (Fan-out/Fan-in pattern)
	fmt.Println("  示例3: Fan-out/Fan-in 模式:")
	input := make(chan int, 10)
	output := make(chan int, 10)

	// Fan-out: 多个 goroutine 读取
	for i := 0; i < 3; i++ {
		go func(id int) {
			for n := range input {
				output <- n * n
			}
		}(i)
	}

	// 发送数据
	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()

	// Fan-in: 收集结果
	go func() {
		wg2 := sync.WaitGroup{}
		wg2.Add(3)
		for i := 0; i < 3; i++ {
			go func() {
				defer wg2.Done()
			}()
		}
		wg2.Wait()
		close(output)
	}()

	time.Sleep(50 * time.Millisecond)
	fmt.Println("  示例3: 扇出扇入完成")
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("    Worker %d 处理任务 %d\n", id, j)
		time.Sleep(10 * time.Millisecond)
		results <- j * 2
	}
}

func init() {
	chapters.Register("awesome", "datapipeline", Run)
}
