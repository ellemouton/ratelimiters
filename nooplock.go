package ratelimiters

import (
	"sync"
	"time"
	"github.com/corverroos/ratelimit"
)

func NewNoopLock(period time.Duration, limit int) *NoopLock {
	return &NoopLock{
		period:  period,
		limit:   limit,
	}
}

type NoopLock struct {
	period time.Duration
	limit  int
	mu      sync.Mutex
}

func (n *NoopLock) Request(resource string) bool {
	n.mu.Lock()
	n.mu.Unlock()

	return true
}

var _ ratelimit.RateLimiter = (*NoopLock)(nil)

