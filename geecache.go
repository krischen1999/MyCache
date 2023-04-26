package GeeCache

import (
	"fmt"
	"log"
	"sync"
)

/*geecache 主体*/

type Group struct {
	name      string
	getter    Getter
	maincache cache
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// 生成group实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	G := &Group{name, getter, cache{
		cacheBytes: cacheBytes,
	}}
	groups[name] = G
	return G
}

// return the named group or nil
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g

}

// get value from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.maincache.get(key); ok {
		log.Printf("[GeeCache] hit")
		return v, nil
	} else {
		log.Printf("[GeeCache] hit failed")
		return g.getLocally(key)
	}

}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

// 从本地获取源数据
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	//将源数据加载到缓存中
	value := ByteView{CloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil

}

func (g *Group) populateCache(key string, value ByteView) {
	g.maincache.update(key, value)
}

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
