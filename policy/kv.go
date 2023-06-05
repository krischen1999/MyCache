package policy

import (
	"container/list"
	"time"
)

type entry struct {
	key     string
	value   Value
	expires *time.Time
}

// ttl
func (ele *entry) expired() (ok bool) {
	if ele.expires == nil {
		ok = true
	} else {
		ok = ele.expires.Before(time.Now())
	}
	return
}

// ttl
func (ele *entry) touch(duration time.Duration) {
	expiration := time.Now().Add(duration)
	ele.expires = &expiration
}

func New(name string, maxBytes int64, onEvicted func(string, Value), ttl time.Duration) Interface {

	if name == "fifo" {
		return NewfifoCache(maxBytes, onEvicted, ttl)
	}
	if name == "lru" {
		return NewLruCache(maxBytes, onEvicted, ttl)
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

func NewfifoCache(maxBytes int64, onEvicted func(string, Value), ttl time.Duration) *fifoCahce {

	return &fifoCahce{
		ttl:       ttl,
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

type Interface interface {
	Get(string) (Value, bool)
	Add(string, Value)
	CleanUp()
}
