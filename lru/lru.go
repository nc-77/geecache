package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	usedBytes int64
	cache     map[string]*list.Element
	ll        *list.List
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	cache := &Cache{
		maxBytes:  maxBytes,
		usedBytes: 0,
		cache:     make(map[string]*list.Element),
		ll:        new(list.List),
		OnEvicted: onEvicted,
	}
	return cache
}

func (c *Cache) Get(key string) (Value, bool) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.ll.MoveToBack(ele)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		c.usedBytes = c.usedBytes - int64(kv.value.Len()) + int64(value.Len())
		kv.value = value
		c.ll.MoveToBack(ele)
		return
	} else {
		ele := c.ll.PushBack(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for c.usedBytes > c.maxBytes {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	if c.ll == nil || c.Len() == 0 {
		return
	}
	frontEle := c.ll.Front()
	if frontEle != nil {
		c.ll.Remove(frontEle)
		kv := frontEle.Value.(*entry)
		delete(c.cache, kv.key)
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}

}

func (c *Cache) Len() int {
	return c.ll.Len()
}
