package utils

import "sync"

type SyncRunner struct {
	tasks chan func()
	wg    sync.WaitGroup
	once  sync.Once
}

// workers: number of worker goroutines
// buffer:  channel buffer size
func NewSyncRunner(workers, buffer int) *SyncRunner {
	r := &SyncRunner{
		tasks: make(chan func(), buffer),
	}

	for range workers {
		go r.worker()
	}

	return r
}

func (r *SyncRunner) worker() {
	for task := range r.tasks {
		task()
		r.wg.Done()
	}
}

// Run schedules a task
func (r *SyncRunner) Run(task func()) {
	r.wg.Add(1)
	r.tasks <- task
}

// Wait blocks until all submitted tasks finish
func (r *SyncRunner) Wait() {
	r.wg.Wait()
}

// Close stops accepting new tasks.
// Call ONLY after all Run() calls are done.
func (r *SyncRunner) Close() {
	r.once.Do(func() {
		close(r.tasks)
	})
}
