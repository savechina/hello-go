package webservice

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore_Add(t *testing.T) {
	tests := []struct {
		name           string
		title          string
		expectedID     int
		expectedTitle  string
		expectedLength int
	}{
		{
			name:           "add first task",
			title:          "Learn Go",
			expectedID:     0,
			expectedTitle:  "Learn Go",
			expectedLength: 1,
		},
		{
			name:           "add second task",
			title:          "Build API",
			expectedID:     1,
			expectedTitle:  "Build API",
			expectedLength: 2,
		},
		{
			name:           "add task with empty title",
			title:          "",
			expectedID:     2,
			expectedTitle:  "",
			expectedLength: 3,
		},
	}

	store := &Store{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := store.Add(tt.title)

			assert.Equal(t, tt.expectedID, task.ID, "task ID should match expected")
			assert.Equal(t, tt.expectedTitle, task.Title, "task title should match expected")
			assert.False(t, task.Completed, "task should not be completed by default")

			tasks := store.List()
			assert.Equal(t, tt.expectedLength, len(tasks), "store should have expected number of tasks")
		})
	}
}

func TestStore_List(t *testing.T) {
	tests := []struct {
		name           string
		setupTasks     []string
		expectedCount  int
		expectedIDs    []int
		expectedTitles []string
	}{
		{
			name:           "empty store",
			setupTasks:     []string{},
			expectedCount:  0,
			expectedIDs:    []int{},
			expectedTitles: []string{},
		},
		{
			name:           "single task",
			setupTasks:     []string{"Task 1"},
			expectedCount:  1,
			expectedIDs:    []int{0},
			expectedTitles: []string{"Task 1"},
		},
		{
			name:           "multiple tasks in order",
			setupTasks:     []string{"Task A", "Task B", "Task C"},
			expectedCount:  3,
			expectedIDs:    []int{0, 1, 2},
			expectedTitles: []string{"Task A", "Task B", "Task C"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &Store{}

			for _, title := range tt.setupTasks {
				store.Add(title)
			}

			tasks := store.List()

			assert.Equal(t, tt.expectedCount, len(tasks), "task count should match expected")

			for i, task := range tasks {
				assert.Equal(t, tt.expectedIDs[i], task.ID, "task ID should match expected order")
				assert.Equal(t, tt.expectedTitles[i], task.Title, "task title should match expected order")
			}
		})
	}
}

func TestStore_ThreadSafety(t *testing.T) {
	store := &Store{}
	numGoroutines := 10
	tasksPerGoroutine := 100
	expectedTotal := numGoroutines * tasksPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < tasksPerGoroutine; j++ {
				title := "task"
				store.Add(title)
			}
		}(i)
	}

	wg.Wait()

	tasks := store.List()
	assert.Equal(t, expectedTotal, len(tasks), "should have correct total number of tasks")

	idSet := make(map[int]bool)
	for _, task := range tasks {
		assert.False(t, idSet[task.ID], "task ID should be unique")
		idSet[task.ID] = true
	}

	assert.Equal(t, expectedTotal, len(idSet), "all task IDs should be unique")
}
