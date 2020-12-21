package ratelimiters_test

import (
	"github.com/corverroos/ratelimit"
	"github.com/ellemouton/ratelimiters"
	"testing"
	"time"
)

func BenchmarkNoop(b *testing.B) {
	ratelimit.Benchmark(b, func() ratelimit.RateLimiter {
		return ratelimiters.NewNoop(time.Millisecond, 10)
	})
}
