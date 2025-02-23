package service

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	m   sync.Map
	t   sync.Map
	ttl time.Duration
}

func NewCache[K comparable, V any](ttl time.Duration) *Cache[K, V] {
	return &Cache[K, V]{
		ttl: ttl,
	}
}

func (c *Cache[K, V]) Load(key K) (value V, ok bool) {
	v, ok := c.m.Load(key)
	if ok {
		created, ok := c.t.Load(key)
		// check if the item's create time + the TTL is still later than now.
		if ok && created.(time.Time).Add(c.ttl).After(time.Now()) {
			// item found, and still within ttl
			return v.(V), ok
		} else {
			// ttl expired, so remove the item from the cache.
			c.Delete(key)
		}
	}
	return value, false
}

func (c *Cache[K, V]) Store(key K, value V) {
	// TODO: implement item limit.
	c.m.Store(key, value)
	c.t.Store(key, time.Now())
}

func (c *Cache[K, V]) Delete(key K) {
	c.m.Delete(key)
	c.t.Delete(key)
}
