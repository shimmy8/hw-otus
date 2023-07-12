package hw04lrucache

type Key string

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
		cachedElem.Value = i
		cache.queue.MoveToFront(cachedElem)
	} else {
		if cache.capacity == cache.queue.Len() {
			cache.queue.Remove(cache.queue.Back())
		}
		elem := cache.queue.PushFront(i)
		cache.items[key] = elem
	}
	return inCache
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cachedElem, inCache := cache.items[key]
	if inCache {
		cache.queue.MoveToFront(cachedElem)
		return cachedElem.Value, inCache
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
