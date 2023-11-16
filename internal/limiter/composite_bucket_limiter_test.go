package limiter_test

import (
	"testing"
	"time"

	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	limitermocks "github.com/petrenko-alex/api-rate-limiter/internal/limiter/mocks"
	"github.com/stretchr/testify/require"
)

func TestCompositeBucketLimiter_SatisfyLimit(t *testing.T) {
	limitStorage := limitermocks.NewMockILimitStorage(t)
	refillRate := limiter.NewRefillRate(3, time.Second*1)

	t.Run("only one bucket, satisfy", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{"login": "lucky"}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("only one bucket, NOT satisfy", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{"login": "looser"}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})

	t.Run("multiple buckets, satisfy all", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{
			"login":    "lucky",
			"ip":       "192.168.1.1",
			"password": "123456",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("multiple buckets, NOT satisfy all", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{
			"login":    "looser",
			"ip":       "192.168.1.2",
			"password": "555",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})

	t.Run("multiple buckets, partial satisfy", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{
			"login":    "lucky",
			"ip":       "192.168.1.1",
			"password": "555",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})

	t.Run("some buckets missed for identity", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{
			"login": "lucky",
			"age":   "18",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.True(t, satisfies)
	})
}

func TestCompositeBucketLimiter_SatisfyLimit_Error(t *testing.T) {
	refillRate := limiter.NewRefillRate(3, time.Second*1)

	t.Run("no limits for identity", func(t *testing.T) {
		key := "age"
		limitStorage := limitermocks.NewMockILimitStorage(t)
		limitStorage.EXPECT().GetLimitsByTypes([]string{key}).Return(&limiter.Limits{}, nil)

		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		satisfies, err := compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{key: "18"},
		)

		require.ErrorIs(t, err, limiter.ErrNoLimitsFound)
		require.False(t, satisfies)
	})

	t.Run("empty identity error", func(t *testing.T) {
		limitStorage := limitermocks.NewMockILimitStorage(t)
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		satisfies, err := compositeLimiter.SatisfyLimit(limiter.UserIdentityDto{})

		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
		require.False(t, satisfies)
	})

	t.Run("invalid limit value", func(t *testing.T) {
		limitType := limiter.LoginLimit
		limits := limiter.Limits{
			limiter.Limit{
				LimitType: limitType,
				Value:     "cast error",
			},
		}
		limitStorage := limitermocks.NewMockILimitStorage(t)
		limitStorage.EXPECT().GetLimitsByTypes([]string{limitType.String()}).Return(&limits, nil)

		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		satisfies, err := compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{limitType.String(): "root"},
		)

		require.ErrorIs(t, err, limiter.ErrInitLimits)
		require.False(t, satisfies)
	})

	t.Run("use another limiter", func(t *testing.T) {
		limitType := limiter.LoginLimit
		limits := limiter.Limits{
			limiter.Limit{
				LimitType: limitType,
				Value:     "3",
			},
		}
		limitStorage := limitermocks.NewMockILimitStorage(t)
		limitStorage.EXPECT().GetLimitsByTypes([]string{limitType.String()}).Return(&limits, nil)

		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		// init limiters with first call
		_, err := compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{limitType.String(): "root"},
		)
		require.NoError(t, err)

		// second call using another identity
		_, err = compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{"age": "18"},
		)
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}
