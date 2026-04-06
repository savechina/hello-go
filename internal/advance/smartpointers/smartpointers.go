package smartpointers

import (
	"fmt"
	"strings"
	"sync"

	"hello/internal/chapters"
)

func init() {
	chapters.Register("advance", "smartpointers", Run)
}

type refCounter struct {
	name     string
	refs     int
	released bool
}

func newRefCounter(name string) *refCounter {
	return &refCounter{name: name, refs: 1}
}

func (r *refCounter) AddRef() int {
	if r == nil || r.released {
		return 0
	}

	r.refs++
	return r.refs
}

func (r *refCounter) Release() int {
	if r == nil || r.refs == 0 {
		return 0
	}

	r.refs--
	if r.refs == 0 {
		r.released = true
	}

	return r.refs
}

func (r *refCounter) Snapshot() string {
	if r == nil {
		return "nil counter"
	}

	return fmt.Sprintf("resource=%s refs=%d released=%t", r.name, r.refs, r.released)
}

type pooledObject struct {
	id      int
	payload []string
}

type objectPool struct {
	pool    sync.Pool
	created int
	nextID  int
}

func newObjectPool() *objectPool {
	op := &objectPool{}
	op.pool.New = func() any {
		op.nextID++
		op.created++
		return &pooledObject{id: op.nextID}
	}
	return op
}

func (o *objectPool) Borrow() *pooledObject {
	return o.pool.Get().(*pooledObject)
}

func (o *objectPool) Return(item *pooledObject) {
	if item == nil {
		return
	}

	item.payload = item.payload[:0]
	o.pool.Put(item)
}

func simulateReferenceCounting(operations []string) string {
	counter := newRefCounter("cache-entry")
	states := []string{counter.Snapshot()}

	for _, op := range operations {
		switch op {
		case "add":
			counter.AddRef()
		case "release":
			counter.Release()
		}
		states = append(states, counter.Snapshot())
	}

	return strings.Join(states, " -> ")
}

func simulatePoolReuse(tasks []string) string {
	pool := newObjectPool()
	used := make([]string, 0, len(tasks))

	for _, task := range tasks {
		item := pool.Borrow()
		item.payload = append(item.payload, task)
		used = append(used, fmt.Sprintf("task=%s object#%d payload=%v", task, item.id, item.payload))
		pool.Return(item)
	}

	return fmt.Sprintf("created=%d %s", pool.created, strings.Join(used, " | "))
}

func processWithCleanup(parts []string) string {
	pool := newObjectPool()
	item := pool.Borrow()
	defer pool.Return(item)

	item.payload = append(item.payload, parts...)
	joined := strings.Join(item.payload, "/")
	clean := len(item.payload)

	return fmt.Sprintf("joined=%s cleanup-items=%d", joined, clean)
}

// Run prints the smart pointer pattern examples.
func Run() {
	fmt.Println("[smartpointers] example 1:", simulateReferenceCounting([]string{"add", "add", "release", "release", "release"}))
	fmt.Println("[smartpointers] example 2:", simulatePoolReuse([]string{"decode", "transform", "flush"}))
	fmt.Println("[smartpointers] example 3:", processWithCleanup([]string{"session", "buffer", "done"}))
}
