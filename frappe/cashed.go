package frappe

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ahmedsat/ebda-cli/config"
)

type cacheEntry struct {
	Doctype  string          `json:"doctype"`
	Name     string          `json:"name"`
	Modified string          `json:"modified"`
	CachedAt time.Time       `json:"cached_at"`
	Data     json.RawMessage `json:"data"`
}

type cacheCall struct {
	done  chan struct{}
	value any
	err   error
}

var frappeCache = struct {
	sync.RWMutex
	items map[string]cacheEntry
	calls map[string]*cacheCall
}{
	items: make(map[string]cacheEntry),
	calls: make(map[string]*cacheCall),
}

// GetCached1 returns a full Frappe document by name, caching the complete
// document returned by Get1. When possible, prefer GetCached1WithModified if a
// list query already gave you the document's modified timestamp.
func GetCached1[T FrappeDoctype](id string) (T, error) {
	return GetCached1WithModified[T](id, "")
}

// GetCached1WithModified returns a full Frappe document by name using the
// caller-provided modified timestamp to validate persistent cache entries.
func GetCached1WithModified[T FrappeDoctype](id, modified string) (T, error) {
	var zero T
	if id == "" {
		return zero, errors.New("id is required")
	}

	doctype := zero.DocTypeName()
	key := cacheKey(doctype, id)

	if doc, ok := memoryCacheGet[T](key, modified); ok {
		return doc, nil
	}

	value, err := doCachedFetch(key, func() (any, error) {
		if doc, ok := memoryCacheGet[T](key, modified); ok {
			return doc, nil
		}

		entry, ok := readDiskCache(doctype, id)
		if ok {
			if modified == "" {
				var err error
				modified, err = remoteModified[T](id)
				if err != nil {
					return zero, err
				}
			}

			if entry.Modified == modified {
				doc, err := decodeCacheEntry[T](entry)
				if err != nil {
					return zero, err
				}
				memoryCacheSet(key, entry)
				return doc, nil
			}
		}

		doc, err := Get1[T](id)
		if err != nil {
			return zero, err
		}

		entry, err = newCacheEntry(doctype, id, doc)
		if err != nil {
			return zero, err
		}

		memoryCacheSet(key, entry)
		_ = writeDiskCache(entry)
		return doc, nil
	})
	if err != nil {
		return zero, err
	}

	doc, ok := value.(T)
	if !ok {
		return zero, fmt.Errorf("cached value has unexpected type %T", value)
	}
	return doc, nil
}

func doCachedFetch(key string, fn func() (any, error)) (any, error) {
	frappeCache.Lock()
	if call, ok := frappeCache.calls[key]; ok {
		frappeCache.Unlock()
		<-call.done
		return call.value, call.err
	}

	call := &cacheCall{done: make(chan struct{})}
	frappeCache.calls[key] = call
	frappeCache.Unlock()

	call.value, call.err = fn()

	frappeCache.Lock()
	delete(frappeCache.calls, key)
	frappeCache.Unlock()
	close(call.done)

	return call.value, call.err
}

func InvalidateCached1[T FrappeDoctype](id string) error {
	var zero T
	if id == "" {
		return errors.New("id is required")
	}

	doctype := zero.DocTypeName()
	key := cacheKey(doctype, id)

	frappeCache.Lock()
	delete(frappeCache.items, key)
	frappeCache.Unlock()

	err := os.Remove(cacheFilePath(doctype, id))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func ClearCache() error {
	frappeCache.Lock()
	frappeCache.items = make(map[string]cacheEntry)
	frappeCache.Unlock()

	return os.RemoveAll(cacheRoot())
}

func memoryCacheGet[T FrappeDoctype](key, modified string) (T, bool) {
	var zero T

	frappeCache.RLock()
	entry, ok := frappeCache.items[key]
	frappeCache.RUnlock()
	if !ok {
		return zero, false
	}
	if modified != "" && entry.Modified != modified {
		return zero, false
	}

	doc, err := decodeCacheEntry[T](entry)
	if err != nil {
		return zero, false
	}
	return doc, true
}

func memoryCacheSet(key string, entry cacheEntry) {
	frappeCache.Lock()
	frappeCache.items[key] = entry
	frappeCache.Unlock()
}

func remoteModified[T FrappeDoctype](id string) (string, error) {
	docs, err := Get[T](
		Filters{NewFilter("name", Eq, id)},
		List{"name", "modified"},
		nil,
	)
	if err != nil {
		return "", err
	}
	if len(docs) == 0 {
		return "", fmt.Errorf("document %q not found", id)
	}

	modified := documentModified(docs[0])
	if modified == "" {
		return "", fmt.Errorf("document %q has empty modified field", id)
	}
	return modified, nil
}

func newCacheEntry[T FrappeDoctype](doctype, id string, doc T) (cacheEntry, error) {
	data, err := json.Marshal(doc)
	if err != nil {
		return cacheEntry{}, err
	}

	name := doc.DocName()
	if name == "" {
		name = id
	}

	return cacheEntry{
		Doctype:  doctype,
		Name:     name,
		Modified: documentModified(doc),
		CachedAt: time.Now(),
		Data:     data,
	}, nil
}

func decodeCacheEntry[T FrappeDoctype](entry cacheEntry) (T, error) {
	var doc T
	err := json.Unmarshal(entry.Data, &doc)
	return doc, err
}

func readDiskCache(doctype, id string) (cacheEntry, bool) {
	data, err := os.ReadFile(cacheFilePath(doctype, id))
	if err != nil {
		return cacheEntry{}, false
	}

	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return cacheEntry{}, false
	}
	if entry.Doctype != doctype || entry.Name != id || len(entry.Data) == 0 {
		return cacheEntry{}, false
	}
	return entry, true
}

func writeDiskCache(entry cacheEntry) error {
	path := cacheFilePath(entry.Doctype, entry.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func cacheKey(doctype, id string) string {
	return doctype + ":" + id
}

func cacheFilePath(doctype, id string) string {
	return filepath.Join(cacheRoot(), safePathPart(doctype), safePathPart(id)+".json")
}

func cacheRoot() string {
	if dir := os.Getenv("EBDA_CLI_FRAPPE_CACHE_DIR"); dir != "" {
		return dir
	}

	root, err := os.UserCacheDir()
	if err != nil {
		root = ".cache"
	}

	hash := sha1.Sum([]byte(config.ErpBaseUrl))
	return filepath.Join(root, "ebda-cli", "frappe", hex.EncodeToString(hash[:])[:12])
}

func safePathPart(s string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(s)
}

func documentModified(doc any) string {
	value := reflect.ValueOf(doc)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ""
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return ""
	}

	field := value.FieldByName("Modified")
	if !field.IsValid() || field.Kind() != reflect.String {
		return ""
	}
	return field.String()
}
