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
	sync.Mutex
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
	lc.Lock()
	defer lc.Unlock()
	val, ok := lc.items[key]
	if ok {
		// если элемент уже был, обновить значение и переместить в начало
		val.Value = value
		lc.queue.MoveToFront(val)
	} else {
		// если не было, добавить
		item := lc.queue.PushFront(value)
		lc.items[key] = item
		// вытолкнуть последний при превышении ёмкости
		if lc.queue.Len() > lc.capacity {
			for k, v := range lc.items {
				if v == lc.queue.Back() {
					delete(lc.items, k)
					break
				}
			}
			lc.queue.Remove(lc.queue.Back())
		}
	}
	return ok
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.Lock()
	defer lc.Unlock()
	val, ok := lc.items[key]
	// если нашелся, переместить в начало
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
