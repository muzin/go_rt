package cache

import (
	"sync"
)

type LRUCache struct {
	cache    map[string]*Entry
	capacity int
	head     *Entry
	tail     *Entry

	mu sync.RWMutex
}

type Entry struct {
	Key   string
	Value interface{}
	pre   *Entry
	next  *Entry
}

func NewLRUCache(cap int) *LRUCache {
	return &LRUCache{
		cache:    make(map[string]*Entry),
		capacity: cap,
	}
}

func (this *LRUCache) Put(key string, val interface{}) interface{} {
	this.mu.Lock()
	defer this.mu.Unlock()

	if existVal, exist := this.cache[key]; exist {
		this.moveToHead(existVal)
		return nil
	}

	e := &Entry{Key: key, Value: val, next: this.head}
	if this.head != nil {
		this.head.pre = e
	}
	this.head = e
	if this.tail == nil {
		this.tail = e
	}
	this.cache[key] = e

	if len(this.cache) <= this.capacity {
		return nil
	}

	removeEntry := this.tail
	this.tail = this.tail.pre
	removeEntry.pre = nil
	this.tail.next = nil
	delete(this.cache, removeEntry.Key)

	return removeEntry.Value
}

func (this *LRUCache) Get(key string) interface{} {
	this.mu.RLock()
	defer this.mu.RUnlock()

	if existVal, exist := this.cache[key]; exist {
		this.moveToHead(existVal)
		return existVal.Value
	}
	return nil
}

// 把元素提到队列头部
func (this *LRUCache) moveToHead(e *Entry) {
	if e == this.head {
		return
	}

	e.pre.next = e.next
	if e == this.tail {
		this.tail = e.pre
	} else {
		e.next.pre = e.pre
	}
	e.pre = nil
	e.next = this.head
	this.head.pre = e
	this.head = e
}
