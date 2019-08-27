package repository

import (
	"github.com/bannovd/evaluator/cache"
	"github.com/bannovd/evaluator/models"

	"strconv"
	"time"
)

// Repository struct
type Repository struct {
	cache *cache.Cache
}

// NewRepository return new repository
func NewRepository(cacheCleanupInterval time.Duration) *Repository {
	c := cache.NewCachedDb(cacheCleanupInterval)

	return &Repository{
		cache: c,
	}
}

// SaveHit func
func (rep *Repository) SaveHit(hit models.Hit) error {
	rep.cache.Add(strconv.FormatInt(time.Now().UnixNano(), 16), hit)
	return nil
}
