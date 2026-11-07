package pokecache

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{key: "https://example.com", val: []byte("testdata")},
		{key: "https://example.com/path", val: []byte("moretestdata")},
	}

	cache := NewCache(interval)
	for _, c := range cases {
		cache.Add(c.key, c.val)
	}

	for _, c := range cases {
		val, ok := cache.Get(c.key)
		if !ok {
			t.Errorf("expected to find key %q", c.key)
			continue
		}
		if string(val) != string(c.val) {
			t.Errorf("key %q: expected %q, got %q", c.key, c.val, val)
		}
	}
}

func TestGetMissing(t *testing.T) {
	cache := NewCache(5 * time.Second)
	val, ok := cache.Get("missing")
	if ok {
		t.Error("expected missing key to return false")
	}
	if val != nil {
		t.Errorf("expected nil value, got %q", val)
	}
}

func TestAddReplaces(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("key", []byte("first"))
	cache.Add("key", []byte("second"))

	val, ok := cache.Get("key")
	if !ok {
		t.Fatal("expected to find key")
	}
	if string(val) != "second" {
		t.Errorf("expected %q, got %q", "second", val)
	}
}

func TestKeysDontOverwrite(t *testing.T) {
	cache := NewCache(5 * time.Second)
	cache.Add("a", []byte("1"))
	cache.Add("b", []byte("2"))

	a, _ := cache.Get("a")
	b, _ := cache.Get("b")
	if string(a) != "1" || string(b) != "2" {
		t.Errorf("keys interfered: a=%q b=%q", a, b)
	}
}

func TestFreshEntryStays(t *testing.T) {
	const interval = 200 * time.Millisecond
	cache := NewCache(interval)
	cache.Add("key", []byte("val"))

	time.Sleep(interval / 4)
	if _, ok := cache.Get("key"); !ok {
		t.Error("expected fresh entry to still be in the cache")
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 50 * time.Millisecond
	cache := NewCache(interval)
	cache.Add("key", []byte("val"))

	time.Sleep(interval * 4)
	if _, ok := cache.Get("key"); ok {
		t.Error("expected old entry to be reaped")
	}
}
