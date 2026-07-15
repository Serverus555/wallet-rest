package cache

import (
	"time"

	"github.com/google/uuid"
	"github.com/maypok86/otter/v2"
)

type InMemory struct {
	cache *otter.Cache[uuid.UUID, int64]
}

func NewInMemory(capacity int, ttl time.Duration) *InMemory {
	return &InMemory{
		cache: otter.Must(&otter.Options[uuid.UUID, int64]{
			MaximumSize:      capacity,
			ExpiryCalculator: otter.ExpiryAccessing[uuid.UUID, int64](ttl),
		}),
	}
}

func (c *InMemory) Put(key uuid.UUID, value int64) {
	c.cache.Set(key, value)
}

func (c *InMemory) Get(key uuid.UUID) (int64, bool) {
	return c.cache.GetIfPresent(key)
}

func (c *InMemory) Delete(key uuid.UUID) {
	c.cache.Invalidate(key)
}
