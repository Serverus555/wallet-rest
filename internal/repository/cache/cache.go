package cache

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Cache interface {
	Put(key uuid.UUID, value int64)
	Get(key uuid.UUID) (int64, bool)
	Delete(key uuid.UUID)
}

type Type string

var (
	MemoryType Type = "memory"
	NoopType   Type = "noop"
	// redis
)

type Config struct {
	CacheType Type          `env:"CACHE_TYPE"`
	Capacity  int           `env:"CACHE_CAPACITY"`
	TTL       time.Duration `env:"CACHE_TTL"`
}

func New(c Config) (Cache, error) {
	switch c.CacheType {
	case MemoryType:
		return NewInMemory(c.Capacity, c.TTL), nil
	case NoopType:
		return NewNoop(), nil
	}
	return nil, fmt.Errorf("invalid cache type: %s", c.CacheType)
}
