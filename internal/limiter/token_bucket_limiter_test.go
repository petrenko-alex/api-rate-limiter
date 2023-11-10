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
		/*refillRate := limiter.RefillRate{
			Time: time.Second * 1,
			Count: 10,
		}*/
		tokenBucketLimiter := limiter.NewTokenBucketLimiter()

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
		refillRate := limiter.RefillRate{
			Time:  time.Second * 1,
			Count: 3,
		}
		tokenBucketLimiter := limiter.NewTokenBucketLimiter()

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
		time.Sleep(refillRate.Time)

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
		/*refillRate := limiter.RefillRate{
			Time: time.Second * 1,
			Count: 10,
		}*/
		tokenBucketLimiter := limiter.NewTokenBucketLimiter()

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
}
