package ratelimiters

import (
	"sync"
	"time"
	"github.com/corverroos/ratelimit"
)

func NewNoop(period time.Duration, limit int) *Noop {
	return &Noop{
		period:  period,
		limit:   limit,
	}
}

type Noop struct {
	period time.Duration
	limit  int
	mu      sync.Mutex
}

func (n *Noop) Request(resource string) bool {
	return true
}

var _ ratelimit.RateLimiter = (*Noop)(nil)
