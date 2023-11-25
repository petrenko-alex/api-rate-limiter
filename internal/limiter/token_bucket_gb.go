package limiter

import (
	"time"
)

// todo: thread safe
type TokenBucketGB struct {
	tokenBucketLimiter ITokenBucketLimitService

	tokenBucketTTL time.Duration
}

func NewTokenBucketGB(tokenBucketLimiter ITokenBucketLimitService, tokenBucketTTL time.Duration) *TokenBucketGB {
	return &TokenBucketGB{
		tokenBucketLimiter: tokenBucketLimiter,
		tokenBucketTTL:     tokenBucketTTL,
	}
}

func (gb *TokenBucketGB) Sweep() {
	buckets := gb.tokenBucketLimiter.GetBuckets()
	if len(buckets) == 0 {
		return
	}

	bucketsToDelete := make([]string, 0)
	for key, bucket := range buckets {

		lastRefill := bucket.GetLastRefill()
		if time.Since(lastRefill) < gb.tokenBucketTTL {
			continue
		}

		bucketsToDelete = append(bucketsToDelete, key)
	}

	if len(bucketsToDelete) == 0 {
		return
	}

	for _, bucketKey := range bucketsToDelete {
		gb.tokenBucketLimiter.SweepBucket(bucketKey)
	}
}
