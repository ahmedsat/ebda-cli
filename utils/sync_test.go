package utils_test

import (
	"bytes"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ahmedsat/ebda-cli/utils"
)

func TestSyncIoWriter(t *testing.T) {
	var buf bytes.Buffer
	w := &utils.SyncIoWriter{Writer: &buf}

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = w.Write([]byte("x"))
		}()
	}
	wg.Wait()

	if buf.Len() != 100 {
		t.Fatalf("buf length = %d, want 100", buf.Len())
	}
}

func TestNewLockableMap(t *testing.T) {
	m := utils.NewLockableMap[string, int]()
	if m.Map == nil {
		t.Fatal("inner map is nil")
	}
	if !m.TryLock() {
		t.Fatal("TryLock should succeed when unlocked")
	}
	m.Unlock()
	if !m.TryRLock() {
		t.Fatal("TryRLock should succeed when unlocked")
	}
	m.RUnlock()
}

func TestLockableMapConcurrentAccess(t *testing.T) {
	m := utils.NewLockableMap[int, int]()
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(2)
		go func(n int) {
			defer wg.Done()
			m.Lock()
			m.Map[n] = n * 2
			m.Unlock()
		}(i)
		go func(n int) {
			defer wg.Done()
			m.RLock()
			_ = m.Map[n]
			m.RUnlock()
		}(i)
	}
	wg.Wait()
}

func TestSyncRunnerAllTasksExecute(t *testing.T) {
	r := utils.NewSyncRunner(4, 16)
	defer r.Close()

	var count atomic.Int64
	for range 20 {
		r.Run(func() error {
			count.Add(1)
			return nil
		})
	}
	if err := r.Wait(); err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if count.Load() != 20 {
		t.Fatalf("executed %d tasks, want 20", count.Load())
	}
}

func TestSyncRunnerWaitReturnsError(t *testing.T) {
	r := utils.NewSyncRunner(2, 8)
	defer r.Close()

	sentinel := errors.New("task failed")
	r.Run(func() error { return sentinel })
	r.Run(func() error { return nil })
	if err := r.Wait(); err == nil {
		t.Fatal("Wait should return an error when a task fails")
	}
}

func TestSyncRunnerConcurrencyRespected(t *testing.T) {
	workers := 3
	r := utils.NewSyncRunner(workers, 20)
	defer r.Close()

	var peak, current atomic.Int64
	for range 15 {
		r.Run(func() error {
			now := current.Add(1)
			for {
				old := peak.Load()
				if now <= old || peak.CompareAndSwap(old, now) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
			current.Add(-1)
			return nil
		})
	}
	if err := r.Wait(); err != nil {
		t.Fatalf("Wait: %v", err)
	}
	if int(peak.Load()) > workers {
		t.Fatalf("peak concurrency = %d, exceeded workers = %d", peak.Load(), workers)
	}
}
