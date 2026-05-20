package frappe

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
)

type cachedTestDoc struct {
	Name     string `json:"name"`
	Modified string `json:"modified"`
	Value    string `json:"value"`
}

func (d cachedTestDoc) DocTypeName() string { return "Cached Test Doc" }
func (d cachedTestDoc) DocName() string     { return d.Name }

func TestCacheEntryRoundTrip(t *testing.T) {
	t.Setenv("EBDA_CLI_FRAPPE_CACHE_DIR", t.TempDir())
	config.ErpBaseUrl = "https://example.test"
	ClearCache()

	doc := cachedTestDoc{
		Name:     "DOC-1",
		Modified: "2026-05-17 10:00:00.000000",
		Value:    "hello",
	}

	entry, err := newCacheEntry(doc.DocTypeName(), doc.DocName(), doc)
	if err != nil {
		t.Fatalf("newCacheEntry: %v", err)
	}
	if err := writeDiskCache(entry); err != nil {
		t.Fatalf("writeDiskCache: %v", err)
	}

	gotEntry, ok := readDiskCache(doc.DocTypeName(), doc.DocName())
	if !ok {
		t.Fatal("readDiskCache did not find written entry")
	}
	if gotEntry.Modified != doc.Modified {
		t.Fatalf("Modified = %q, want %q", gotEntry.Modified, doc.Modified)
	}

	got, err := decodeCacheEntry[cachedTestDoc](gotEntry)
	if err != nil {
		t.Fatalf("decodeCacheEntry: %v", err)
	}
	if got != doc {
		t.Fatalf("decoded doc = %#v, want %#v", got, doc)
	}
}

func TestMemoryCacheGetHonorsModified(t *testing.T) {
	key := cacheKey("Cached Test Doc", "DOC-2")
	doc := cachedTestDoc{Name: "DOC-2", Modified: "new", Value: "cached"}
	entry, err := newCacheEntry(doc.DocTypeName(), doc.DocName(), doc)
	if err != nil {
		t.Fatalf("newCacheEntry: %v", err)
	}
	memoryCacheSet(key, entry)
	t.Cleanup(func() {
		frappeCache.Lock()
		delete(frappeCache.items, key)
		frappeCache.Unlock()
	})

	if _, ok := memoryCacheGet[cachedTestDoc](key, "old"); ok {
		t.Fatal("memory cache returned stale modified entry")
	}

	got, ok := memoryCacheGet[cachedTestDoc](key, "new")
	if !ok {
		t.Fatal("memory cache did not return matching modified entry")
	}
	if got != doc {
		t.Fatalf("got %#v, want %#v", got, doc)
	}
}

func TestDoCachedFetchSharesConcurrentFetch(t *testing.T) {
	key := "shared-fetch"
	var calls atomic.Int64
	var wg sync.WaitGroup
	start := make(chan struct{})

	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			value, err := doCachedFetch(key, func() (any, error) {
				calls.Add(1)
				time.Sleep(20 * time.Millisecond)
				return "ok", nil
			})
			if err != nil {
				t.Errorf("doCachedFetch: %v", err)
				return
			}
			if value != "ok" {
				t.Errorf("value = %v, want ok", value)
			}
		}()
	}

	close(start)
	wg.Wait()

	if calls.Load() != 1 {
		t.Fatalf("fetch calls = %d, want 1", calls.Load())
	}
}
