package limiter

import (
	"time"
)

// todo: thread safe
type TokenBucketGB struct {
	tokenBucketLimiter *CompositeBucketLimiter // todo: интерфейс для token bucket limiter'ов

	tokenBucketTTL time.Duration
}

func NewTokenBucketGB(tokenBucketLimiter *CompositeBucketLimiter, tokenBucketTTL time.Duration) *TokenBucketGB {
	return &TokenBucketGB{
		tokenBucketLimiter: tokenBucketLimiter,
		tokenBucketTTL:     tokenBucketTTL,
	}
}

func (gb *TokenBucketGB) Sweep() {
	limiters := gb.tokenBucketLimiter.GetLimiters()
	if len(limiters) == 0 {
		return
	}

	for _, limiter := range limiters {
		buckets := limiter.GetBuckets() // todo: method to interface?
		if len(buckets) == 0 {
			continue
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
			continue
		}

		for _, bucketKey := range bucketsToDelete {
			delete(buckets, bucketKey)
		}
	}
}
