package limiter_test

import (
	"testing"
	"time"

	"api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestTokenBucketLimiter_SatisfyLimit(t *testing.T) {
	identity := limiter.UserIdentityDto{"ip": "192.168.1.1"}

	t.Run("limit reached", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(3, time.Second*1)
		tokenBucketLimiter := limiter.NewTokenBucketLimiter("ip", bucketSize, refillRate)

		// 3 requests allowed
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th and following requests denied
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.False(t, satisfies)
			require.NoError(t, err)
		}

	})

	t.Run("simple refill", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(3, time.Second*1)
		tokenBucketLimiter := limiter.NewTokenBucketLimiter("ip", bucketSize, refillRate)

		// 3 requests allowed
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th denied
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		// after refill rate time we can make requests again
		time.Sleep(refillRate.GetTime())

		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th denied
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("partial refill", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(10, time.Second*1)
		tokenBucketLimiter := limiter.NewTokenBucketLimiter("ip", bucketSize, refillRate)

		// 3 requests allowed
		for i := 0; i < bucketSize; i++ {
			satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
			require.True(t, satisfies)
			require.NoError(t, err)
		}

		// 4th denied
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		// refill rate = 10req/sec,
		// so we can make one more request after 0.1sec
		time.Sleep(time.Millisecond * 100)

		// one more request is allowed
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// followings not
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("dynamic request cost", func(t *testing.T) {
		refillRate := limiter.NewRefillRate(10, time.Second*1)
		tokenBucketLimiter := limiter.NewTokenBucketLimiter("ip", 10, refillRate)

		// request with cost of 6 tokens after some time
		time.Sleep(time.Millisecond * 300)
		tokenBucketLimiter.SetRequestCost(6)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// check allowed requests remained (0 because only 4 tokens left)
		require.Zero(t, tokenBucketLimiter.GetRequestsAllowed())

		// wait 200 ms to refill and make request with cost of 5 tokens
		time.Sleep(time.Millisecond * 200)
		tokenBucketLimiter.SetRequestCost(5)
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// check allowed requests remained (0 because only 1 token left)
		require.Zero(t, tokenBucketLimiter.GetRequestsAllowed())

		// wait for full refill and check requests allowed (2 request allowed with cost of 5 each)
		time.Sleep(time.Second * 1)
		require.Equal(t, 2, tokenBucketLimiter.GetRequestsAllowed())
	})
}
