package limiter

import "time"

type TokenBucket struct {
	bucketSize uint
	refillRate RefillRate

	tokensCount int
}

type RefillRate struct {
	Count uint
	Time  time.Duration
}
