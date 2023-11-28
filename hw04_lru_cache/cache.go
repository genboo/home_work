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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	val, ok := lc.items[key]
	if ok {
		val.Value = value
		lc.queue.MoveToFront(val)
	} else {
		item := lc.queue.PushFront(value)
		lc.items[key] = item
		if lc.queue.Len() > lc.capacity {
			lc.queue.Remove(lc.queue.Back())
		}
	}
	return ok
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	val, ok := lc.items[key]
	if ok {
		lc.queue.MoveToFront(val)
		return val.Value, true
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
