package limiter

import (
	"time"
)

type TokenBucket struct { // todo: thread safe
	size       int
	refillRate RefillRate

	tokensCount int
	lastRefill  time.Time
}

type RefillRate struct {
	count int
	time  time.Duration // todo: what if not 1 sec
}

func NewTokenBucket(size int, refillRate RefillRate) TokenBucket {
	return TokenBucket{
		size:       size,
		refillRate: refillRate,

		tokensCount: size,
		lastRefill:  time.Now(),
	}
}

func NewRefillRate(count int, time time.Duration) RefillRate {
	return RefillRate{
		count: count,
		time:  time,
	}
}

func (b *TokenBucket) Refill() {
	const nsInSec = 1e9
	timePassed := time.Since(b.lastRefill)
	tokensToAdd := int64(timePassed) * int64(b.refillRate.count) / nsInSec

	b.tokensCount = min(b.tokensCount+int(tokensToAdd), b.size)
	b.lastRefill = time.Now()
}

func (b *TokenBucket) GetTokenCount() int {
	return b.tokensCount
}

func (b *TokenBucket) GetToken(tokenCount int) {
	b.tokensCount -= tokenCount
}

func (r RefillRate) GetTime() time.Duration {
	return r.time
}

func (r RefillRate) GetCount() int {
	return r.count
}
