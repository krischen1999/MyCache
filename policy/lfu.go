package policy

import (
	"container/heap"
	"time"
)

type PriorityQueue struct {
	ttl       time.Duration
	nbytes    int64
	maxBytes  int64
	cache     map[string]*lfuEntry
	pq        *priorityqueue
	OnEvicted func(key string, value Value)
}

func (p PriorityQueue) Get(key string) (value Value, ok bool) {
	if e, ok := p.cache[key]; ok {
		e.referenced()
		heap.Fix(p.pq, e.index)
		return e.entry.value, ok
	}
	return
}

func (p PriorityQueue) Add(key string, value Value) {
	if e, ok := p.cache[key]; ok {
		//更新value
		p.nbytes += int64(value.Len()) - int64(e.entry.value.Len())
		e.entry.value = value
		e.referenced()
		heap.Fix(p.pq, e.index)
	} else {

		e := &lfuEntry{0, entry{key, value, nil}, 0}
		heap.Push(p.pq, e)
		p.cache[key] = e
		p.nbytes += int64(len(key)) + int64(value.Len())
	}
	for p.maxBytes != 0 && p.maxBytes < p.nbytes {
		p.Remove()
	}
}

func (p PriorityQueue) CleanUp() {
	for _, e := range *p.pq {
		if e.entry.expired(p.ttl) {

			kv := heap.Remove(p.pq, e.index).(*lfuEntry).entry
			delete(p.cache, kv.key)
			p.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
			if p.OnEvicted != nil {
				p.OnEvicted(kv.key, kv.value)
			}
		}
	}
}

func (p PriorityQueue) Remove() {
	e := p.pq.Pop().(*lfuEntry)
	delete(p.cache, e.entry.key)
	p.nbytes -= int64(len(e.entry.key)) + int64(e.entry.value.Len())
	if p.OnEvicted != nil {
		p.OnEvicted(e.entry.key, e.entry.value)
	}
}
