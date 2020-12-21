package ratelimiters_test

import (
	"github.com/corverroos/ratelimit"
	limiters "github.com/ellemouton/ratelimiters"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestChannelWindow(t *testing.T) {
	l := limiters.NewChannelWindow(time.Hour, 10)
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			require.True(t, l.Request(""))
			wg.Done()
		}()
	}
	wg.Wait()
	require.False(t, l.Request(""))
}

func BenchmarkChannelWindow(b *testing.B) {
	ratelimit.Benchmark(b, func() ratelimit.RateLimiter {
		return limiters.NewChannelWindow(time.Millisecond, 10)
	})
}