package api

import (
	"sync"
	"time"
)

type messageLimiter struct {
	mu        sync.Mutex
	tokens    float64
	rate      float64
	burst     float64
	last      time.Time
	violation int
}

func newMessageLimiter(rate, burst float64) *messageLimiter {
	return &messageLimiter{tokens: burst, rate: rate, burst: burst, last: time.Now()}
}

func (l *messageLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.tokens += now.Sub(l.last).Seconds() * l.rate
	if l.tokens > l.burst {
		l.tokens = l.burst
	}
	l.last = now

	if l.tokens < 1 {
		l.violation++

		return false
	}

	l.tokens--
	l.violation = 0

	return true
}

func (l *messageLimiter) ExceededRepeatedly() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.violation >= messageRateLimitCloses
}
