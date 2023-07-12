package hw04lrucache

type Key string

type CacheItem struct {
	value interface{}
	key   Key
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (cache *lruCache) Set(key Key, i interface{}) bool {
	cachedElem, inCache := cache.items[key]

	if inCache {
		cachedElem.Value = CacheItem{value: i, key: key}
		cache.queue.MoveToFront(cachedElem)
	} else {
		if cache.capacity == cache.queue.Len() {
			queueBack := cache.queue.Back()
			cacheKey := queueBack.Value.(CacheItem).key
			cache.queue.Remove(queueBack)
			delete(cache.items, cacheKey)
		}
		elem := cache.queue.PushFront(CacheItem{value: i, key: key})
		cache.items[key] = elem
	}
	return inCache
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cachedElem, inCache := cache.items[key]
	if inCache {
		cache.queue.MoveToFront(cachedElem)
		return cachedElem.Value.(CacheItem).value, inCache
	}
	return nil, inCache
}

func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
