package lrucache

import (
	"container/list"
	"sync"
)

type entry[V any] struct {
	key   string
	value V
}

// Cache — потокобезопасный LRU-кэш с дженериками
type Cache[V any] struct {
	mu       sync.RWMutex
	capacity int
	list     *list.List
	items    map[string]*list.Element
}

func New[V any](capacity int) *Cache[V] {
	if capacity < 0 {
		capacity = 0
	}
	return &Cache[V]{
		capacity: capacity,
		list:     list.New(),
		items:    make(map[string]*list.Element, capacity),
	}
}

// Set добавляет или обновляет элемент в кэше.
// Если элемент уже есть — перемещает его в начало.
// Если кэш переполнен — удаляет наименее используемый (последний) элемент.
func (c *Cache[V]) Set(key string, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity == 0 {
		return
	}

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*entry[V]).value = value
		return
	}

	if c.list.Len() >= c.capacity {
		last := c.list.Back()
		if last != nil {
			c.list.Remove(last)
			delete(c.items, last.Value.(*entry[V]).key)
		}
	}

	e := &entry[V]{key: key, value: value}
	elem := c.list.PushFront(e)
	c.items[key] = elem
}

// Get возвращает элемент по ключу.
// Если найден — перемещает в начало и возвращает (value, true).
// Если не найден — возвращает (zero-value, false).
func (c *Cache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry[V]).value, true
	}

	var zero V
	return zero, false
}

func (c *Cache[V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.Init()
	c.items = make(map[string]*list.Element, c.capacity)
}
