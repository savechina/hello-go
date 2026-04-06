package concurrency

import (
	"fmt"
	"sync"
	"time"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("basic", "concurrency", Run)
}

func main() {
	Run()
}

// Run executes the concurrency chapter examples.
func Run() {
	examples := []string{
		fmt.Sprintf("1) goroutine + channel => %v", generateSequence(4)),
		fmt.Sprintf("2) worker pattern + WaitGroup => %v", squareJobs([]int{2, 3, 4}, 2)),
		"3) select timeout => " + waitForSignal(2*time.Millisecond, 20*time.Millisecond),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

type job struct {
	index int
	value int
}

type result struct {
	index int
	value int
}

func generateSequence(count int) []int {
	values := make([]int, 0, count)
	stream := make(chan int)

	go func() {
		defer close(stream)
		for i := 1; i <= count; i++ {
			stream <- i
		}
	}()

	for value := range stream {
		values = append(values, value)
	}

	return values
}

func squareJobs(values []int, workers int) []int {
	jobs := make(chan job, len(values))
	results := make(chan result, len(values))

	var wg sync.WaitGroup
	for workerID := 0; workerID < workers; workerID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jobItem := range jobs {
				results <- result{index: jobItem.index, value: jobItem.value * jobItem.value}
			}
		}()
	}

	for index, value := range values {
		jobs <- job{index: index, value: value}
	}
	close(jobs)

	wg.Wait()
	close(results)

	ordered := make([]int, len(values))
	for item := range results {
		ordered[item.index] = item.value
	}

	return ordered
}

func waitForSignal(produceDelay time.Duration, timeout time.Duration) string {
	signal := make(chan string, 1)

	go func() {
		time.Sleep(produceDelay)
		signal <- "work finished"
	}()

	select {
	case message := <-signal:
		return message
	case <-time.After(timeout):
		return "timeout reached"
	}
}
