package cache

import (
	"sync"
	"time"
)

// Cache struct cache
type Cache struct {
	sync.RWMutex
	items           map[string]Item
	cleanupInterval time.Duration
}

// Item struct cache item
type Item struct {
	Value   interface{}
	Created time.Time
}

// NewCache Initializing a new memory cache
func NewCache(cleanupInterval time.Duration) *Cache {
	items := make(map[string]Item)

	cache := Cache{
		items:           items,
		cleanupInterval: cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// Add a cache item
func (c *Cache) Add(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = Item{
		Value:   value,
		Created: time.Now(),
	}
}

// StartGC func
func (c *Cache) StartGC() {
	go c.GC()
}

// GC Garbage Collection
func (c *Cache) GC() {
	for {
		<-time.After(c.cleanupInterval)
		if c.items == nil {
			return
		}
		c.storeAndRemoveItems()
	}
}

// storeAndRemoveItems store and removes all items
func (c *Cache) storeAndRemoveItems() {
	c.Lock()
	defer c.Unlock()

	for k, _ := range c.items {
		//Todo: create func for store 'k' using any provider (file, db, external api service, etc.)
		delete(c.items, k)
	}
}
