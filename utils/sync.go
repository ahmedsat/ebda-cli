package utils

import (
	"io"
	"sync"
)

type SyncRunner struct {
	tasks chan func() error
	wg    sync.WaitGroup
	once  sync.Once
	err   error
}

// workers: number of worker goroutines
// buffer:  channel buffer size
func NewSyncRunner(workers, buffer int) *SyncRunner {
	r := &SyncRunner{
		tasks: make(chan func() error, buffer),
	}

	for range workers {
		go r.worker()
	}

	return r
}

func (r *SyncRunner) worker() {
	for task := range r.tasks {
		if r.err != nil {
			r.wg.Done()
			continue
		}
		err := task()
		if err != nil {
			r.err = err
		}
		r.wg.Done()
	}
}

// Run schedules a task
func (r *SyncRunner) Run(task func() error) {
	if r.err != nil {
		return
	}
	r.wg.Add(1)
	r.tasks <- task
}

// Wait blocks until all submitted tasks finish
func (r *SyncRunner) Wait() error {
	r.wg.Wait()
	return r.err
}

// Close stops accepting new tasks.
// Call ONLY after all Run() calls are done.
func (r *SyncRunner) Close() {
	r.once.Do(func() {
		close(r.tasks)
	})
}

type SyncIoWriter struct {
	sync.Mutex
	io.Writer
}

func (w *SyncIoWriter) Write(b []byte) (int, error) {
	w.Lock()
	defer w.Unlock()
	return w.Writer.Write(b)
}

type LockableMap[K comparable, V any] struct {
	sync.RWMutex
	Map map[K]V
}

func NewLockableMap[K comparable, V any]() LockableMap[K, V] {
	return LockableMap[K, V]{
		Map: make(map[K]V),
	}
}

func (m *LockableMap[K, V]) Lock() {
	m.RWMutex.Lock()
}

func (m *LockableMap[K, V]) Unlock() {
	m.RWMutex.Unlock()
}

func (m *LockableMap[K, V]) TryLock() bool {
	return m.RWMutex.TryLock()
}

func (m *LockableMap[K, V]) RLock() {
	m.RWMutex.RLock()
}

func (m *LockableMap[K, V]) RUnlock() {
	m.RWMutex.RUnlock()
}

func (m *LockableMap[K, V]) TryRLock() bool {
	return m.RWMutex.TryRLock()
}
