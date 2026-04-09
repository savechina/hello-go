package datapipeline

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	tests := []struct {
		name       string
		workerID   int
		jobs       []int
		wantResult []int
	}{
		{
			name:       "single job",
			workerID:   1,
			jobs:       []int{5},
			wantResult: []int{10},
		},
		{
			name:       "multiple jobs",
			workerID:   2,
			jobs:       []int{1, 2, 3, 4, 5},
			wantResult: []int{2, 4, 6, 8, 10},
		},
		{
			name:       "zero value",
			workerID:   3,
			jobs:       []int{0},
			wantResult: []int{0},
		},
		{
			name:       "negative values",
			workerID:   4,
			jobs:       []int{-5, -3, -1},
			wantResult: []int{-10, -6, -2},
		},
		{
			name:       "empty jobs channel",
			workerID:   5,
			jobs:       []int{},
			wantResult: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobs := make(chan int, len(tt.jobs))
			results := make(chan int, len(tt.jobs))
			var wg sync.WaitGroup

			for _, j := range tt.jobs {
				jobs <- j
			}
			close(jobs)

			wg.Add(1)
			go worker(tt.workerID, jobs, results, &wg)

			go func() {
				wg.Wait()
				close(results)
			}()

			var gotResults []int
			for r := range results {
				gotResults = append(gotResults, r)
			}

			assert.Equal(t, len(tt.wantResult), len(gotResults), "result count mismatch")
			for _, expected := range tt.wantResult {
				assert.Contains(t, gotResults, expected, "expected result not found")
			}
		})
	}
}

func TestWorkerCallsWaitGroupDone(t *testing.T) {
	jobs := make(chan int, 1)
	results := make(chan int, 1)
	var wg sync.WaitGroup

	jobs <- 42
	close(jobs)

	wg.Add(1)
	worker(1, jobs, results, &wg)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("worker did not call wg.Done()")
	}

	select {
	case r := <-results:
		assert.Equal(t, 84, r, "expected 42 * 2 = 84")
	default:
		t.Fatal("no result received from worker")
	}
}

func TestWorkerPool(t *testing.T) {
	tests := []struct {
		name        string
		numWorkers  int
		numJobs     int
		jobValues   []int
		wantResults []int
	}{
		{
			name:        "3 workers 5 jobs",
			numWorkers:  3,
			numJobs:     5,
			jobValues:   []int{1, 2, 3, 4, 5},
			wantResults: []int{2, 4, 6, 8, 10},
		},
		{
			name:        "1 worker 1 job",
			numWorkers:  1,
			numJobs:     1,
			jobValues:   []int{100},
			wantResults: []int{200},
		},
		{
			name:        "2 workers 10 jobs",
			numWorkers:  2,
			numJobs:     10,
			jobValues:   []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			wantResults: []int{20, 40, 60, 80, 100, 120, 140, 160, 180, 200},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobs := make(chan int, tt.numJobs)
			results := make(chan int, tt.numJobs)
			var wg sync.WaitGroup

			for w := 1; w <= tt.numWorkers; w++ {
				wg.Add(1)
				go worker(w, jobs, results, &wg)
			}

			for _, j := range tt.jobValues {
				jobs <- j
			}
			close(jobs)

			go func() {
				wg.Wait()
				close(results)
			}()

			var gotResults []int
			for r := range results {
				gotResults = append(gotResults, r)
			}

			assert.Equal(t, len(tt.wantResults), len(gotResults), "result count mismatch")
			for _, expected := range tt.wantResults {
				assert.Contains(t, gotResults, expected, "expected result not found")
			}
		})
	}
}

func TestGracefulShutdown(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		wantDone bool
	}{
		{
			name:     "completes before timeout",
			timeout:  200 * time.Millisecond,
			wantDone: true,
		},
		{
			name:     "timeout occurs",
			timeout:  10 * time.Millisecond,
			wantDone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})

			go func() {
				time.Sleep(50 * time.Millisecond)
				close(done)
			}()

			select {
			case <-done:
				if !tt.wantDone {
					t.Fatal("expected timeout but task completed")
				}
			case <-time.After(tt.timeout):
				if tt.wantDone {
					t.Fatal("expected task completion but timed out")
				}
			}
		})
	}
}

func TestPipelineIntegration(t *testing.T) {
	jobs := make(chan int, 5)
	results := make(chan int, 5)
	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	expectedResults := []int{2, 4, 6, 8, 10}
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	var gotResults []int
	for r := range results {
		gotResults = append(gotResults, r)
	}

	assert.Equal(t, len(expectedResults), len(gotResults), "should receive all 5 results")
	for _, expected := range expectedResults {
		assert.Contains(t, gotResults, expected, "should contain expected doubled value")
	}
}
