package cfapi

import (
	"errors"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	if cache == nil {
		t.Fatal("NewCache() returned nil")
	}
	if cache.ttl != 5*time.Minute {
		t.Errorf("ttl = %v, want %v", cache.ttl, 5*time.Minute)
	}
	if cache.entries == nil {
		t.Error("entries map should be initialized")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	// Set a value
	cache.Set("key1", "value1")

	// Get the value
	value, ok := cache.Get("key1")
	if !ok {
		t.Error("Get() should return true for existing key")
	}
	if value != "value1" {
		t.Errorf("Get() = %v, want %v", value, "value1")
	}
}

func TestCache_Get_NonExistent(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	value, ok := cache.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for non-existent key")
	}
	if value != nil {
		t.Errorf("Get() should return nil for non-existent key, got %v", value)
	}
}

func TestCache_Get_Expired(t *testing.T) {
	// Create cache with very short TTL
	cache := NewCache(1 * time.Millisecond)

	cache.Set("key", "value")

	// Wait for expiration
	time.Sleep(5 * time.Millisecond)

	value, ok := cache.Get("key")
	if ok {
		t.Error("Get() should return false for expired key")
	}
	if value != nil {
		t.Errorf("Get() should return nil for expired key, got %v", value)
	}
}

func TestCache_SetWithTTL(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	// Set with custom short TTL
	cache.SetWithTTL("key", "value", 1*time.Millisecond)

	// Should exist initially
	value, ok := cache.Get("key")
	if !ok {
		t.Error("Get() should return true immediately after SetWithTTL")
	}
	if value != "value" {
		t.Errorf("Get() = %v, want %v", value, "value")
	}

	// Wait for expiration
	time.Sleep(5 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("key")
	if ok {
		t.Error("Get() should return false after TTL expires")
	}
}

func TestCache_Delete(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	cache.Set("key", "value")

	// Verify it exists
	_, ok := cache.Get("key")
	if !ok {
		t.Error("Key should exist before delete")
	}

	// Delete
	cache.Delete("key")

	// Verify it's gone
	_, ok = cache.Get("key")
	if ok {
		t.Error("Key should not exist after delete")
	}
}

func TestCache_Delete_NonExistent(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	// Should not panic
	cache.Delete("nonexistent")
}

func TestCache_Clear(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if cache.Size() != 3 {
		t.Errorf("Size() = %v, want 3", cache.Size())
	}

	cache.Clear()

	if cache.Size() != 0 {
		t.Errorf("Size() after Clear() = %v, want 0", cache.Size())
	}
}

func TestCache_Size(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	if cache.Size() != 0 {
		t.Errorf("Size() of new cache = %v, want 0", cache.Size())
	}

	cache.Set("key1", "value1")
	if cache.Size() != 1 {
		t.Errorf("Size() = %v, want 1", cache.Size())
	}

	cache.Set("key2", "value2")
	if cache.Size() != 2 {
		t.Errorf("Size() = %v, want 2", cache.Size())
	}

	cache.Delete("key1")
	if cache.Size() != 1 {
		t.Errorf("Size() after delete = %v, want 1", cache.Size())
	}
}

func TestCache_GetOrSet_CacheHit(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	// Pre-populate cache
	cache.Set("key", "cached_value")

	fnCalled := false
	value, err := cache.GetOrSet("key", func() (interface{}, error) {
		fnCalled = true
		return "new_value", nil
	})

	if err != nil {
		t.Errorf("GetOrSet() error = %v", err)
	}
	if fnCalled {
		t.Error("Function should not be called on cache hit")
	}
	if value != "cached_value" {
		t.Errorf("GetOrSet() = %v, want %v", value, "cached_value")
	}
}

func TestCache_GetOrSet_CacheMiss(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	fnCalled := false
	value, err := cache.GetOrSet("key", func() (interface{}, error) {
		fnCalled = true
		return "new_value", nil
	})

	if err != nil {
		t.Errorf("GetOrSet() error = %v", err)
	}
	if !fnCalled {
		t.Error("Function should be called on cache miss")
	}
	if value != "new_value" {
		t.Errorf("GetOrSet() = %v, want %v", value, "new_value")
	}

	// Value should now be cached
	cached, ok := cache.Get("key")
	if !ok {
		t.Error("Value should be cached after GetOrSet")
	}
	if cached != "new_value" {
		t.Errorf("Cached value = %v, want %v", cached, "new_value")
	}
}

func TestCache_GetOrSet_FunctionError(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	expectedErr := errors.New("fetch error")
	value, err := cache.GetOrSet("key", func() (interface{}, error) {
		return nil, expectedErr
	})

	if err != expectedErr {
		t.Errorf("GetOrSet() error = %v, want %v", err, expectedErr)
	}
	if value != nil {
		t.Errorf("GetOrSet() value = %v, want nil", value)
	}

	// Value should not be cached on error
	_, ok := cache.Get("key")
	if ok {
		t.Error("Value should not be cached after error")
	}
}

func TestCache_ComplexTypes(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	// Store a struct
	problem := Problem{ContestID: 1325, Index: "A", Name: "Test"}
	cache.Set("problem", problem)

	value, ok := cache.Get("problem")
	if !ok {
		t.Error("Get() should return true for stored struct")
	}

	retrieved, ok := value.(Problem)
	if !ok {
		t.Error("Retrieved value should be Problem type")
	}
	if retrieved.ContestID != 1325 {
		t.Errorf("retrieved.ContestID = %v, want %v", retrieved.ContestID, 1325)
	}
}

func TestCache_Overwrite(t *testing.T) {
	cache := NewCache(5 * time.Minute)

	cache.Set("key", "value1")
	cache.Set("key", "value2")

	value, ok := cache.Get("key")
	if !ok {
		t.Error("Get() should return true")
	}
	if value != "value2" {
		t.Errorf("Get() = %v, want value2 (overwritten)", value)
	}
	if cache.Size() != 1 {
		t.Errorf("Size() = %v, want 1 (no duplicate)", cache.Size())
	}
}

func TestCacheEntry_Fields(t *testing.T) {
	entry := CacheEntry{
		Value:      "test",
		Expiration: time.Now().Add(5 * time.Minute),
	}

	if entry.Value != "test" {
		t.Errorf("Value = %v, want test", entry.Value)
	}
	if entry.Expiration.Before(time.Now()) {
		t.Error("Expiration should be in the future")
	}
}
