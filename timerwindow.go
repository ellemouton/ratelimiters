package ratelimiters

import (
	"sync"
	"time"
	"github.com/corverroos/ratelimit"
)

func NewTimerWindow(period time.Duration, limit int) *TimerWindow {
	counts := make(map[string]int)

	timer := &TimerWindow{
		period:  period,
		limit:   limit,
		nowFunc: time.Now,
		counts: counts,
	}

	go timer.resetForever()

	return  timer
}

func (t *TimerWindow) resetForever() {
	ticker := time.NewTicker(t.period)

	for {
		select{
		case <- ticker.C:
			t.mu.Lock()
			t.counts = make(map[string]int)
			t.mu.Unlock()
		}
	}
}

type TimerWindow struct {
	period time.Duration
	limit  int

	counts  map[string]int
	mu sync.Mutex

	nowFunc func() time.Time
}

func (t *TimerWindow) Request(resource string) bool {
	t.mu.Lock()
	//defer t.mu.Unlock() // `defer` has non-trivial overhead
	c := t.counts[resource]
	c++
	t.counts[resource] = c
	ret := c <= t.limit

	t.mu.Unlock() // remove if using `defer`

	return ret
}


var _ ratelimit.RateLimiter = (*TimerWindow)(nil)
