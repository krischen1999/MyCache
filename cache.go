package MyCache

/*并发控制*/
import (
	"GeeCache/policy"
	"sync"
	"time"
)

// 该文件为lru添加并发特性
type cache struct {
	ttl        time.Duration
	mu         sync.RWMutex //互斥锁
	plo        Interface
	cacheBytes int64
}

type Interface interface {
	Get(string) (policy.Value, bool)
	Add(string, policy.Value)
	CleanUp()
}

func (c *cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//延迟初始化
	if c.plo == nil {
		c.plo = policy.NewLruCache(c.cacheBytes, nil, c.ttl)

	}

	c.plo.Add(key, value)

}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.plo == nil {
		return
	}
	if v, ok := c.plo.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

func (c *cache) statrClenupTimer() {
	duration := c.ttl
	//ClenUp触发时间不少于一分钟
	if duration < time.Second {
		duration = time.Second
	}
	ticker := time.Tick(duration)
	go func() {
		for {
			select {
			case <-ticker:
				c.mu.Lock()
				c.plo.CleanUp()
				c.mu.Unlock()
			}
		}
	}()

}
