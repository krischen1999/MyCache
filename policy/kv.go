package policy

import (
	"container/list"
	"time"
)

type entry struct {
	key      string
	value    Value
	updateAt *time.Time
}

// ttl
func (ele *entry) expired(duration time.Duration) (ok bool) {
	if ele.updateAt == nil {
		ok = true
	} else {
		ok = ele.updateAt.Add(duration).Before(time.Now())
	}
	return
}

// ttl
func (ele *entry) touch() {
	//ele.updateAt=time.Now()
	nowTime := time.Now()
	ele.updateAt = &nowTime
}

func New(name string, maxBytes int64, onEvicted func(string, Value), ttl time.Duration) Interface {

	if name == "fifo" {
		return NewFifoCache(maxBytes, onEvicted, ttl)
	}
	if name == "lru" {
		return NewLruCache(maxBytes, onEvicted, ttl)
	}
	if name == "lfu" {
		return NewLfuCache(maxBytes, onEvicted, ttl)
	}

	return nil
}

func NewLruCache(maxBytes int64, onEvicted func(string, Value), ttl time.Duration) *lru {

	return &lru{
		ttl:       ttl,
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func NewFifoCache(maxBytes int64, onEvicted func(string, Value), ttl time.Duration) *fifoCahce {

	return &fifoCahce{
		ttl:       ttl,
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func NewLfuCache(maxBytes int64, onEvicted func(string, Value), ttl time.Duration) *lfuCache {
	queue := priorityqueue(make([]*lfuEntry, 0))
	return &lfuCache{
		ttl:       ttl,
		maxBytes:  maxBytes,
		pq:        &queue,
		cache:     make(map[string]*lfuEntry),
		OnEvicted: onEvicted,
	}
}

type Interface interface {
	Get(string) (Value, bool)
	Add(string, Value)
	CleanUp()
}
