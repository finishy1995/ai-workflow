package data

import (
	"github.com/patrickmn/go-cache"
)

type Cache[T any] struct {
	cache *cache.Cache
}

func (c *Cache[T]) Get(k string) (T, bool) {
	v, ok := c.cache.Get(k)
	if !ok {
		var result T
		return result, false
	}
	return v.(T), ok
}

func (c *Cache[T]) Set(k string, v T) {
	c.cache.Set(k, v, cache.NoExpiration)
}

func (c *Cache[T]) Items() map[string]T {
	tmp := make(map[string]T)
	for k, item := range c.cache.Items() {
		tmp[k] = item.Object.(T)
	}
	return tmp
}

func (c *Cache[T]) Delete(k string) {
	c.cache.Delete(k)
}

func (c *Cache[T]) Count() int {
	return c.cache.ItemCount()
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}
