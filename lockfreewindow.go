package ratelimiters

import (
	"github.com/corverroos/ratelimit"
	"runtime"
	"sync/atomic"
	"time"
)

// LockFreeWindow is not 100% accurate since only 1024 buckets are used for resources uses a bad string hashing func
// (other hash functions are too slow).
// The lock-free observation pattern was inspired by:
// https://github.com/prometheus/client_golang/blob/master/prometheus/histogram.go

const numBuckets = 1024

func NewLockFreeWindow(period time.Duration, limit int) LockFreeWindow {
	l := LockFreeWindow{
		period:     period,
		limit:      limit,
		nowFunc:    time.Now,
		counts:     [2]*data{{}, {}},
		numBuckets: numBuckets,
	}

	l.counts[0].buckets = make([]uint64, numBuckets)
	l.counts[1].buckets = make([]uint64, numBuckets)

	go l.ReadForever()

	return l
}

type data struct {
	timestamp int64
	buckets   []uint64
	count     uint64
}

type LockFreeWindow struct {
	period time.Duration
	limit  int

	numBuckets     int
	counts         [2]*data
	countAndHotIdx uint64

	nowFunc func() time.Time
}

func (w *LockFreeWindow) Request(resource string) bool {
	return w.Write(resource)
}

func (w *LockFreeWindow) Write(resource string) bool {
	n := atomic.AddUint64(&w.countAndHotIdx, 1)
	hotCounts := w.counts[n>>63]
	index := w.getIndex(resource)
	atomic.AddUint64(&hotCounts.buckets[index], 1)
	hc := atomic.LoadUint64(&hotCounts.buckets[index])
	atomic.AddUint64(&hotCounts.count, 1)

	return int(hc) <= w.limit
}

func (w *LockFreeWindow) ReadForever() {
	ticker := time.NewTicker(w.period)

	for {
		select {
		case <-ticker.C:
			n := atomic.AddUint64(&w.countAndHotIdx, 1<<63)
			count := n & ((1 << 63) - 1)
			hotCounts := w.counts[n>>63]
			coldCounts := w.counts[(^n)>>63]

			for count != atomic.LoadUint64(&coldCounts.count) {
				runtime.Gosched()
			}

			atomic.AddUint64(&hotCounts.count, atomic.LoadUint64(&coldCounts.count))
			atomic.StoreUint64(&coldCounts.count, 0)
			for i := 0; i < w.numBuckets; i++ {
				atomic.StoreUint64(&coldCounts.buckets[i], 0)
			}
		}
	}
}

func (w *LockFreeWindow) getIndex(resource string) int {
	res := 0
	for i := range []rune(resource) {
		res += i
	}
	return res % w.numBuckets
}

var _ ratelimit.RateLimiter = (*LockFreeWindow)(nil)
