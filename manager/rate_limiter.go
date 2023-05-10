package manager

import (
	"loadbalancer/APIRequest"
	"math"
	"sync"
	"time"
)

// RateLimiterIntf is to be implemented by all rate limiting algorithms
type RateLimiterIntf interface {
	IsRequestAllowed(req *APIRequest.APIRequest) bool
}

// TokenBucket implements RateLimiterIntf interface
type TokenBucket struct {
	rate            int64
	maxTokens       int64
	currTokens      int64
	lastestRefillTs time.Time
	lock            sync.Mutex
}

func NewTokenBucket(rate, maxTokens int64) *TokenBucket {
	return &TokenBucket{
		rate:            rate,
		maxTokens:       maxTokens,
		currTokens:      0,
		lastestRefillTs: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now()
	end := time.Since(tb.lastestRefillTs)
	tokensToBeAdded := (end.Nanoseconds() * tb.rate) / 1000000000
	tb.currTokens = int64(math.Min(float64(tb.currTokens+tokensToBeAdded),
		float64(tb.maxTokens)))
	tb.lastestRefillTs = now
}

func (tb *TokenBucket) IsRequestAllowed(req *APIRequest.APIRequest) bool {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	tb.refill()
	if tb.currTokens >= 1 {
		tb.currTokens -= 1
		return true
	}

	return false
}
