package limiter_test

import (
	"testing"

	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestCompositeBucketLimiter_SatisfyLimit(t *testing.T) {
	t.Run("only one bucket, satisfy", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter()

		identity := limiter.UserIdentityDto{"login": "lucky"}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("only one bucket, NOT satisfy", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter()

		identity := limiter.UserIdentityDto{"login": "looser"}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})

	t.Run("multiple buckets, satisfy all", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter()

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
		compositeLimiter := limiter.NewCompositeBucketLimiter()

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
		compositeLimiter := limiter.NewCompositeBucketLimiter()

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
		compositeLimiter := limiter.NewCompositeBucketLimiter()

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
	t.Run("no buckets for identity", func(t *testing.T) {
		compositeLimiter := limiter.NewCompositeBucketLimiter()

		identity := limiter.UserIdentityDto{
			"age": "18",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.ErrorIs(t, err, limiter.ErrNoRulesFound)
		require.False(t, satisfies)
	})
}
