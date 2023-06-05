package policy

import (
	"time"
)

type PriorityQueue struct {
	ttl      time.Duration
	nbytes   int64
	maxBytes int64

	OnEvicted func(key string, value Value)
}
