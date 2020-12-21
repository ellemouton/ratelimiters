package ratelimiters

import (
	"github.com/corverroos/ratelimit"
	"sync"
	"time"
)

func NewChannelWindow(period time.Duration, limit int) *ChannelWindow {
	return &ChannelWindow{
		period:  period,
		limit:   limit,
		nowFunc: time.Now,
	}
}

type ChannelWindow struct {
	period time.Duration
	limit  int

	current time.Time
	counts  map[string]chan struct{}

	mu sync.Mutex

	nowFunc func() time.Time
}

func (b *ChannelWindow) Request(resource string) bool {
	b.mu.Lock()
	//defer b.mu.Unlock() // `defer` has non-trivial overhead

	bucket := b.nowFunc().Truncate(b.period)
	if bucket == b.current {
		select {
		case b.counts[resource] <- struct{}{}:
			b.mu.Unlock() // remove if using `defer`
			return true
		default:
			b.mu.Unlock() // remove if using `defer`
			return false
		}
	}

	b.current = bucket
	b.counts = make(map[string]chan struct{})
	b.counts[resource] = make(chan struct{}, b.limit - 1)

	b.mu.Unlock() // remove if using `defer`
	return true
}

var _ ratelimit.RateLimiter = (*ChannelWindow)(nil)