package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	sync.RWMutex

	size       int
	refillRate RefillRate

	tokensCount int
	lastRefill  time.Time
}

type RefillRate struct {
	count int
	time  time.Duration
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
	b.Lock()

	const nsInSec = 1e9
	timePassed := time.Since(b.lastRefill)
	tokensToAdd := int64(timePassed) * int64(b.refillRate.count) / int64(nsInSec*b.refillRate.time.Seconds())

	b.tokensCount = min(b.tokensCount+int(tokensToAdd), b.size)
	if tokensToAdd > 0 {
		b.lastRefill = time.Now()
	}

	b.Unlock()
}

func (b *TokenBucket) GetTokenCount() int {
	b.RLock()
	defer b.RUnlock()

	return b.tokensCount
}

func (b *TokenBucket) GetToken(tokenCount int) {
	b.Lock()

	b.tokensCount -= tokenCount

	b.Unlock()
}

func (b *TokenBucket) Reset() {
	b.Lock()

	b.tokensCount = b.size
	b.lastRefill = time.Now()

	b.Unlock()
}

func (r RefillRate) GetTime() time.Duration {
	return r.time
}

func (r RefillRate) GetCount() int {
	return r.count
}
