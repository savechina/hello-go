package chapters

import (
	"fmt"
	"sort"
	"sync"
)

type Runner func()

var (
	mu sync.RWMutex

	registry = map[string]map[string]Runner{
		"basic":    {},
		"advance":  {},
		"awesome":  {},
		"algo":     {},
		"leetcode": {},
		"quiz":     {},
	}
)

func Register(level string, name string, runner Runner) {
	mu.Lock()
	defer mu.Unlock()

	chaptersByLevel, ok := registry[level]
	if !ok {
		panic(fmt.Sprintf("unknown chapter level: %s", level))
	}

	if runner == nil {
		panic(fmt.Sprintf("nil runner for %s %s", level, name))
	}

	if _, exists := chaptersByLevel[name]; exists {
		panic(fmt.Sprintf("chapter already registered: %s %s", level, name))
	}

	chaptersByLevel[name] = runner
}

func Lookup(level string, name string) (Runner, bool) {
	mu.RLock()
	defer mu.RUnlock()

	chaptersByLevel, ok := registry[level]
	if !ok {
		return nil, false
	}

	runner, ok := chaptersByLevel[name]
	return runner, ok
}

func Names(level string) []string {
	mu.RLock()
	defer mu.RUnlock()

	chaptersByLevel, ok := registry[level]
	if !ok {
		return nil
	}

	names := make([]string, 0, len(chaptersByLevel))
	for name := range chaptersByLevel {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}
