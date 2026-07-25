package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		val      []byte
		expected []byte
	}{
		{
			name:     "simple string",
			key:      "test-key",
			val:      []byte("hello world"),
			expected: []byte("hello world"),
		},
		{
			name:     "empty bytes",
			key:      "empty",
			val:      []byte{},
			expected: []byte{},
		},
		{
			name:     "json bytes",
			key:      "api-response",
			val:      []byte(`{"name":"pikachu","id":25}`),
			expected: []byte(`{"name":"pikachu","id":25}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewCache(5 * time.Second)

			cache.Add(tt.key, tt.val)
			val, ok := cache.Get(tt.key)

			if !ok {
				t.Errorf("Get(%q) expected ok=true, got false", tt.key)
				return
			}
			if string(val) != string(tt.expected) {
				t.Errorf("Get(%q) = %q; expected %q", tt.key, val, tt.expected)
			}
		})
	}
}

func TestCacheGetMissing(t *testing.T) {
	cache := NewCache(5 * time.Second)

	val, ok := cache.Get("non-existent")

	if ok {
		t.Errorf("Get missing key expected ok=false, got true")
	}
	if val != nil {
		t.Errorf("Get missing key expected nil, got %v", val)
	}
}

func TestCacheOverwrite(t *testing.T) {
	cache := NewCache(5 * time.Second)

	cache.Add("key", []byte("first"))
	cache.Add("key", []byte("second"))

	val, ok := cache.Get("key")
	if !ok {
		t.Errorf("Get after overwrite expected ok=true")
		return
	}
	if string(val) != "second" {
		t.Errorf("Overwrite failed: got %q, expected %q", val, "second")
	}
}

func TestCacheReapLoop(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)

	// Add "old" entry first
	cache.Add("old", []byte("old"))

	// Wait longer than interval so "old" is eligible for reaping
	time.Sleep(150 * time.Millisecond)

	// Add "fresh" entry after the wait
	cache.Add("fresh", []byte("new"))

	// Wait a bit for reap loop to run (but less than interval for fresh)
	time.Sleep(20 * time.Millisecond)

	// Fresh should still exist (added 20ms ago, interval=50ms)
	_, ok := cache.Get("fresh")
	if !ok {
		t.Errorf("Fresh entry should not be reaped")
	}

	// Old should be reaped (added >150ms ago, interval=50ms)
	_, ok = cache.Get("old")
	if ok {
		t.Errorf("Old entry should have been reaped")
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Add("key", []byte("value"))
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("key")
		}
		done <- true
	}()

	// Wait for both
	<-done
	<-done

	// Should not panic or deadlock
	val, ok := cache.Get("key")
	if !ok || string(val) != "value" {
		t.Errorf("Concurrent access failed")
	}
}