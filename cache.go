package MyCache

/*并发控制*/
import (
	"MyCache/policy"
	"sync"
	"time"
)

// 该文件为cache 添加并发特性
// 并为cache添加ttl功能
type cache struct {
	ttl        time.Duration
	mu         sync.RWMutex //互斥锁
	polName    string
	pol        Interface
	cacheBytes int64
}

type Interface interface {
	Get(string) (policy.Value, *time.Time, bool)
	Add(string, policy.Value)
	CleanUp(ttl time.Duration)
	Len() int
}

func (c *cache) Add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//延迟初始化
	if c.pol == nil {
		c.pol = policy.New(c.polName, c.cacheBytes, nil)
	}
	c.pol.Add(key, value)

}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.pol == nil {
		return
	}
	if v, t, ok := c.pol.Get(key); ok && t.Add(c.ttl).After(time.Now()) {
		return v.(ByteView), ok
	}
	return
}

func (c *cache) startClenUpTimer() {
	duration := c.ttl
	//TTL不小于一秒钟
	if duration < time.Second {
		duration = time.Second
	}
	ticker := time.Tick(duration)
	//ClenUp触发间隔跟ttl一致
	go func() {
		for {
			//select {
			//case <-ticker:
			//	c.mu.Lock()
			//	c.pol.CleanUp(duration)
			//	c.mu.Unlock()
			//}
			<-ticker
			c.mu.Lock()
			c.pol.CleanUp(duration)
			c.mu.Unlock()
		}
	}()

}
