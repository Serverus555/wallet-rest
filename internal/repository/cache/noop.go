package cache

import (
	"github.com/google/uuid"
)

type Noop struct {
}

func NewNoop() *Noop {
	return &Noop{}
}

func (c *Noop) Put(key uuid.UUID, value int64) {
}

func (c *Noop) Get(key uuid.UUID) (int64, bool) {
	return 0, false
}

func (c *Noop) Delete(key uuid.UUID) {
}
