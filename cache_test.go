package MyCache

import (
	"fmt"
	"testing"
	"time"
)

func Test_cache_statrClenupTimer(t *testing.T) {
	c := &cache{
		cacheBytes: 100,
		polName:    "lru",
		ttl:        time.Second,
	}

	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	c.Add(k1, ByteView{data: []byte(v1)})
	c.Add(k2, ByteView{data: []byte(v2)})
	c.Add(k3, ByteView{data: []byte(v3)})
	c.startClenUpTimer()
	for c.pol.Len() != 0 {
		time.Sleep(500 * time.Microsecond)
		fmt.Println(c.pol.Len())
	}

}
