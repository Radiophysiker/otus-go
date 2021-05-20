package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex sync.Mutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	newItem := cacheItem{
		key:   key,
		value: value,
	}

	cachedValue, wasInCache := cache.items[key]
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if wasInCache {
		cache.queue.MoveToFront(cachedValue)
		cachedValue.Value = newItem
		return true
	}
	cache.items[key] = cache.queue.PushFront(newItem)
	if cache.queue.Len() > cache.capacity {
		last := cache.queue.Back()
		cache.queue.Remove(last)
		delete(cache.items, last.Value.(cacheItem).key)
	}
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cachedValue, wasInCache := cache.items[key]
	if !wasInCache {
		return nil, false
	}
	cache.queue.MoveToFront(cachedValue)
	return cachedValue.Value.(cacheItem).value, true
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}
