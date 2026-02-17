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
