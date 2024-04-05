package limiter_test

import (
	"testing"
	"time"

	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestTokenBucketGB_Sweep(t *testing.T) {
	size := 3
	ttl := time.Millisecond * 100
	refillRate := limiter.NewRefillRate(1, time.Hour*1) // "disable" auto refill with long rate
	tokenBucketLimiter := getTokenBucketLimiter(t, size, refillRate)
	gb := limiter.NewTokenBucketGB(tokenBucketLimiter, ttl)

	// init token buckets
	tokenBucketLimiter.SetRequestCost(size)
	identities := []limiter.UserIdentityDto{
		{limiter.IPLimit.String(): "192.168.1.1"},
		{limiter.LoginLimit.String(): "root"},
		{limiter.LoginLimit.String(): "alex"},
		{limiter.PasswordLimit.String(): "123456"},
	}
	for _, identity := range identities {
		tokenBucketLimiter.SatisfyLimit(identity)
		tokenBucketLimiter.ResetLimit(identity)
	}
	require.Len(t, tokenBucketLimiter.GetBuckets(), len(identities))

	// drain one identity to check GB sweeping only full buckets (despite ttl expired)
	tokenBucketLimiter.SatisfyLimit(identities[2])

	// sleep to make buckets refill date outdated
	time.Sleep(ttl)

	// run db
	gb.Sweep()

	require.Len(t, tokenBucketLimiter.GetBuckets(), 1)
}

func getTokenBucketLimiter(t *testing.T, size int, refillRate limiter.RefillRate) *limiter.CompositeBucketLimiter {
	t.Helper()

	return limiter.NewCompositeBucketLimiter(
		getMockLimitStorage(
			t,
			[]limiter.LimitType{limiter.LoginLimit, limiter.IPLimit, limiter.PasswordLimit},
			[]int{size, size, size},
		),
		refillRate,
	)
}
