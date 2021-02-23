package hw04_lru_cache // nolint:golint,stylecheck
import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		cap:        capacity,
		queue:      NewList(),
		cacheMutex: sync.Mutex{},
		cache:      make(map[Key]*ListItem),
	}
}

type lruCache struct {
	cap        int
	queue      List
	cacheMutex sync.Mutex
	cache      map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.cacheMutex.Lock()
	defer l.cacheMutex.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.queue.MoveToFront(elem)
		elem.Value = cacheItem{
			key:   key,
			value: value,
		}
		return true
	}

	elem := cacheItem{
		key:   key,
		value: value,
	}
	newelem := l.queue.PushFront(elem)
	l.cache[key] = newelem

	if l.queue.Len() > l.cap {
		last := l.queue.Back()
		l.queue.Remove(last)

		cacheItem, ok := last.Value.(cacheItem)
		if !ok {
			return false
		}

		delete(l.cache, cacheItem.key)
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.cacheMutex.Lock()
	defer l.cacheMutex.Unlock()

	if elem, ok := l.cache[key]; ok {
		l.queue.MoveToFront(elem)
		cacheItem, ok := elem.Value.(cacheItem)
		if !ok {
			return nil, false
		}
		return cacheItem.value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.cacheMutex.Lock()
	defer l.cacheMutex.Unlock()

	l.queue = NewList()
	l.cache = make(map[Key]*ListItem)
}
