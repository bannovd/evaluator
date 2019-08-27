package repository

import (
	"github.com/bannovd/evaluator/cache"
	"github.com/bannovd/evaluator/models"
	"strconv"
	"time"
)

// Repository struct
type Repository struct {
	timeout time.Duration
	cache   *cache.Cache
}

// NewRepository return new repository
func NewRepository() *Repository {
	c := cache.NewCache(5 * time.Second)

	return &Repository{
		timeout: time.Second,
		cache:   c,
	}
}

// SaveHit func
func (rep *Repository) SaveHit(hit models.Hit) error {
	rep.cache.Add(strconv.FormatInt(time.Now().UnixNano(), 16), hit)
	return nil
}
