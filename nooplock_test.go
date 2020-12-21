package ratelimiters_test

import (
	"github.com/corverroos/ratelimit"
	limiter "github.com/ellemouton/ratelimiters"
	"testing"
	"time"
)

func BenchmarkNoopLock(b *testing.B) {
	ratelimit.Benchmark(b, func() ratelimit.RateLimiter {
		return limiter.NewNoopLock(time.Millisecond, 10)
	})
}