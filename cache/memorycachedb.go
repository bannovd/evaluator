package cache

import (
	"fmt"
	"github.com/bannovd/evaluator/models"
	"os"
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

// NewCachedDb Initializing a new memory cache
func NewCachedDb(cleanupInterval time.Duration) *Cache {
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

	path := "data.csv"

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Print("File Does Not Exist: ")
		}
		fmt.Println(err)
	}
	defer file.Close()

	for k, v := range c.items {
		hit := v.Value.(models.Hit)
		_, err = file.WriteString(fmt.Sprintf("%s;%s\n", hit.Type, hit.Value))
		delete(c.items, k)
	}

	_ = file.Sync()
}
