package limiter_test

import (
	"testing"
	"time"

	"api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestTokenBucketLimiter_SatisfyLimit(t *testing.T) {
	bucketKey := "ip"
	identity := limiter.UserIdentityDto{bucketKey: "192.168.1.1"}

	t.Run("limit reached", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(3, time.Second*1)
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, bucketSize, refillRate)

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
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, bucketSize, refillRate)

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
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, bucketSize, refillRate)

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
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, 10, refillRate)

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

	t.Run("tricky refill rate #1", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(3, time.Second*3) // same as 1t/1sec
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, bucketSize, refillRate)

		tokenBucketLimiter.SetRequestCost(3)
		tokenBucketLimiter.SatisfyLimit(identity) // waste all tokens
		tokenBucketLimiter.SetRequestCost(1)

		// expect 1 token after 1 sec
		time.Sleep(time.Second * 1)
		require.Equal(t, 1, tokenBucketLimiter.GetRequestsAllowed())
	})

	t.Run("tricky refill rate #2", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(125, time.Second*150) // 125t/2.5min = same as 0.8(3)t/1sec
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, bucketSize, refillRate)

		tokenBucketLimiter.SetRequestCost(3)
		tokenBucketLimiter.SatisfyLimit(identity) // waste all tokens
		tokenBucketLimiter.SetRequestCost(1)

		// expect 1 full token after 2 sec
		time.Sleep(time.Second * 1)
		require.Equal(t, 0, tokenBucketLimiter.GetRequestsAllowed())

		time.Sleep(time.Second * 1)
		require.Equal(t, 1, tokenBucketLimiter.GetRequestsAllowed())
	})

	t.Run("multiple identity", func(t *testing.T) {
		bucketSize := 3
		refillRate := limiter.NewRefillRate(3, time.Second*1)
		tokenBucketLimiter := limiter.NewTokenBucketLimiter(bucketKey, bucketSize, refillRate)

		// waste all tokens for first ip
		identity1 := limiter.UserIdentityDto{bucketKey: "192.168.1.1"}
		tokenBucketLimiter.SetRequestCost(3)
		satisfies, err := tokenBucketLimiter.SatisfyLimit(identity1)
		require.True(t, satisfies)
		require.NoError(t, err)

		// check no more allowed for first ip
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity1)
		require.False(t, satisfies)
		require.NoError(t, err)

		// check allowed for another ip
		identity2 := limiter.UserIdentityDto{bucketKey: "192.155.10.32"}
		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity2)
		require.True(t, satisfies)
		require.NoError(t, err)

		// wait to refill, try one more
		time.Sleep(time.Second * 1)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity1)
		require.True(t, satisfies)
		require.NoError(t, err)

		satisfies, err = tokenBucketLimiter.SatisfyLimit(identity2)
		require.True(t, satisfies)
		require.NoError(t, err)
	})
}
