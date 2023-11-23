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

	// fill with some token buckets and drain it
	tokenBucketLimiter.SetRequestCost(size)
	identities := []limiter.UserIdentityDto{
		{limiter.IPLimit.String(): "192.168.1.1"},
		{limiter.LoginLimit.String(): "root"},
		{limiter.LoginLimit.String(): "alex"},
		{limiter.PasswordLimit.String(): "123456"},
	}
	for _, identity := range identities {
		tokenBucketLimiter.SatisfyLimit(identity)
	}

	// sleep to make buckets refill date outdated
	time.Sleep(ttl)

	// manually reset limits for some identities, to make its refill date fresh
	tokenBucketLimiter.ResetLimit(identities[2])
	tokenBucketLimiter.ResetLimit(identities[3])

	// drain identities with fresh refill date
	tokenBucketLimiter.SatisfyLimit(identities[2])
	tokenBucketLimiter.SatisfyLimit(identities[3])

	// run db
	gb.Sweep()

	// check first identity pair satisfies
	for _, identity := range identities[0:2] {
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.True(t, satisfies)
	}

	// check second pair not
	for _, identity := range identities[2:4] {
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.NoError(t, err)
		require.False(t, satisfies)
	}
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
