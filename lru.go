package q

import (
	"container/list"
)

// Pair is the value of a doubly linked list.
type Pair struct {
	key   int
	value int
}

// LRUCache is a Least recently used cache implemented upon doubly linked list and hash map.
type LRUCache struct {
	cap int                   // capacity
	l   *list.List            // doubly linked list
	m   map[int]*list.Element // hash table for value existence check
}

// NewLRUCache initializes a new LRUCache.
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		cap: capacity,
		l:   new(list.List),
		m:   make(map[int]*list.Element, capacity),
	}
}

// Get returns a existing value by key.
func (c *LRUCache) Get(key int) int {
	if node, ok := c.m[key]; ok {
		val := node.Value.(Pair).value
		c.l.MoveToFront(node)
		return val
	}
	return -1
}

// Put inserts or updates value by key.
func (c *LRUCache) Put(key int, value int) {
	if node, ok := c.m[key]; ok {
		c.l.MoveToFront(node)
		node.Value = Pair{key: key, value: value}
	} else {
		if c.l.Len() == c.cap {
			idx := c.l.Back().Value.(Pair).key
			delete(c.m, idx)
			c.l.Remove(c.l.Back())
		}

		node := Pair{
			key:   key,
			value: value,
		}

		ptr := c.l.PushFront(node)
		c.m[key] = ptr
	}
}
