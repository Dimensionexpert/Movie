package cache

import (
	"container/list"

	"github.com/Dimensionexpert/movieLibrary/internal/movie"
)

type entry struct {
	key   int
	value movie.Movie
}

type Cache struct {
	capacity int
	item     map[int]*list.Element
	order    *list.List
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		item:     make(map[int]*list.Element),
		order:    list.New(),
	}
}

func (c *Cache) Get(key int) (movie.Movie, bool) {
	elem, found := c.item[key]
	if !found {
		return movie.Movie{}, false
	}

	c.order.MoveToFront(elem)
	e := elem.Value.(entry)
	return e.value, true

}

func (c *Cache) Put(key int, value movie.Movie) {
	// key exist - update the value, mark recently used
	if elem, found := c.item[key]; found {
		elem.Value = entry{key: key, value: value}
		c.order.MoveToFront(elem)
		return
	}
	// new key; checking capacity beforehand
	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			oldestEntry := oldest.Value.(entry)
			delete(c.item, oldestEntry.key)
		}

	}

	// adding new item at front
	elem := c.order.PushFront(entry{key: key, value: value})
	c.item[key] = elem
}
