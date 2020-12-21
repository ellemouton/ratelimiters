package ratelimiters

import (
	"sync"
	"time"
	"github.com/corverroos/ratelimit"
)

// Coffee rate limiter is a WIP.

func NewCoffee(period time.Duration, limit int) *Coffee{
	return &Coffee{
		period:  period,
		limit:   limit,
		nowFunc: time.Now,
		mm: newMapMutex(),
		counts: map[string]burst{},
	}
}

type burst struct {
	t time.Time
	count int
}

type Coffee struct {
	period time.Duration
	limit  int

	counts map[string]burst
	mm *mapMutex

	nowFunc func() time.Time
}

type mapMutex struct {
	locks map[string]chan struct{}
	mu sync.Mutex
}

func newMapMutex() *mapMutex{
	return &mapMutex{
		locks: make(map[string]chan struct{}),
	}
}

func (m *mapMutex) lock(res string) {
	m.mu.Lock()
	if _, ok := m.locks[res]; !ok {
		m.locks[res] = make(chan struct{}, 1)
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	select {
	case <- m.locks[res]:
	}
}

func (m *mapMutex) unlock(res string) {
	m.locks[res] <- struct{}{}
}

func (c *Coffee) Request(resource string) bool {
	c.mm.lock(resource)
	defer c.mm.unlock(resource)

	i, ok := c.counts[resource]
	if !ok || i.t != c.nowFunc().Truncate(c.period){
		c.counts[resource] = burst{
			t:     c.nowFunc().Truncate(c.period),
			count: 1,
		}
		return true
	}

	c.counts[resource] = burst{
		t:     i.t,
		count: i.count+1,
	}
	return i.count+1 <= c.limit
}

var _ ratelimit.RateLimiter = (*Coffee)(nil)
