package cache

import (
	"testing"
)

func TestLRUCache_Get(t *testing.T) {
	t.Run("tt.name", func(t *testing.T) {

		lruCache := NewLRUCache(10)

		lruCache.Put("1", 1)
		head := lruCache.head
		if head.Key != "1" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("2", 2)
		head = lruCache.head
		if head.Key != "2" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("3", 3)
		head = lruCache.head
		if head.Key != "3" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("4", 4)
		head = lruCache.head
		if head.Key != "4" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("5", 5)
		head = lruCache.head
		if head.Key != "5" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("6", 6)
		head = lruCache.head
		if head.Key != "6" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("7", 7)
		head = lruCache.head
		if head.Key != "7" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("8", 8)
		head = lruCache.head
		if head.Key != "8" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("9", 9)
		head = lruCache.head
		if head.Key != "9" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("10", 10)
		head = lruCache.head
		if head.Key != "10" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("11", 11)
		head = lruCache.head
		if head.Key != "11" {
			t.Errorf("head.Key != %s", head.Key)
		}

		lruCache.Put("4", 4)
		head = lruCache.head
		if head.Key != "4" {
			t.Errorf("head.Key != %s", head.Key)
		}

		t.Logf("")

	})
}
