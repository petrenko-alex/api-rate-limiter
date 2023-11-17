package limiter_test

import (
	"testing"
	"time"

	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	limitermocks "github.com/petrenko-alex/api-rate-limiter/internal/limiter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCompositeBucketLimiter_SatisfyLimit(t *testing.T) {
	refillRate := limiter.NewRefillRate(3, time.Second*1)

	t.Run("only one bucket, satisfy", func(t *testing.T) {
		usedType := limiter.LoginLimit
		types := []limiter.LimitType{usedType}
		limitStorage := getMockLimitStorage(t, types, []string{"3"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{usedType.String(): "lucky"}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("only one bucket, NOT satisfy", func(t *testing.T) {
		usedType := limiter.LoginLimit
		types := []limiter.LimitType{usedType}
		limitStorage := getMockLimitStorage(t, types, []string{"3"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)
		compositeLimiter.SetRequestCost(4)

		identity := limiter.UserIdentityDto{usedType.String(): "looser"}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})

	t.Run("multiple buckets, satisfy all", func(t *testing.T) {
		types := []limiter.LimitType{
			limiter.LoginLimit,
			limiter.IPLimit,
			limiter.PasswordLimit,
		}
		limitStorage := getMockLimitStorage(t, types, []string{"3", "3", "3"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{
			types[0].String(): "lucky",
			types[1].String(): "192.168.1.1",
			types[2].String(): "123456",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("multiple buckets, satisfy none", func(t *testing.T) {
		types := []limiter.LimitType{
			limiter.LoginLimit,
			limiter.IPLimit,
			limiter.PasswordLimit,
		}
		limitStorage := getMockLimitStorage(t, types, []string{"3", "3", "3"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)
		compositeLimiter.SetRequestCost(4)

		identity := limiter.UserIdentityDto{
			types[0].String(): "looser",
			types[1].String(): "192.168.1.2",
			types[2].String(): "555",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})

	t.Run("multiple buckets, partial satisfy", func(t *testing.T) {
		types := []limiter.LimitType{
			limiter.LoginLimit,
			limiter.IPLimit,
			limiter.PasswordLimit,
		}
		limitStorage := getMockLimitStorage(t, types, []string{"3", "3", "2"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)
		compositeLimiter.SetRequestCost(3)

		identity := limiter.UserIdentityDto{
			types[0].String(): "lucky",
			types[1].String(): "192.168.1.1",
			types[2].String(): "555",
		}
		satisfies, err := compositeLimiter.SatisfyLimit(identity)

		require.NoError(t, err)
		require.False(t, satisfies)
	})
}

func TestCompositeBucketLimiter_SatisfyLimit_Error(t *testing.T) {
	refillRate := limiter.NewRefillRate(3, time.Second*1)

	t.Run("no limits for identity", func(t *testing.T) {
		key := "age"
		limitStorage := getMockLimitStorage(t, []limiter.LimitType{}, []string{})
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
		usedType := limiter.LoginLimit
		types := []limiter.LimitType{usedType}
		limitStorage := getMockLimitStorage(t, types, []string{"cast error"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		satisfies, err := compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{usedType.String(): "root"},
		)

		require.ErrorIs(t, err, limiter.ErrInitLimits)
		require.False(t, satisfies)
	})

	t.Run("use another limiter", func(t *testing.T) {
		usedType := limiter.LoginLimit
		types := []limiter.LimitType{usedType}
		limitStorage := getMockLimitStorage(t, types, []string{"3"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		// init limiters with first call
		_, err := compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{usedType.String(): "root"},
		)
		require.NoError(t, err)

		// second call using another identity
		_, err = compositeLimiter.SatisfyLimit(
			limiter.UserIdentityDto{"age": "18"},
		)
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})

	t.Run("some buckets missed for identity", func(t *testing.T) {
		usedType := limiter.LoginLimit
		types := []limiter.LimitType{usedType}
		limitStorage := getMockLimitStorage(t, types, []string{"3"})
		compositeLimiter := limiter.NewCompositeBucketLimiter(limitStorage, refillRate)

		identity := limiter.UserIdentityDto{
			usedType.String(): "lucky",
			"age":             "18",
		}
		_, err := compositeLimiter.SatisfyLimit(identity)

		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}

func getMockLimitStorage(t *testing.T, types []limiter.LimitType, values []string) *limitermocks.MockILimitStorage {
	t.Helper()

	mockLimits := make(limiter.Limits, 0, len(types))
	for i, limitType := range types {
		mockLimits = append(mockLimits, limiter.Limit{
			LimitType: limitType,
			Value:     values[i],
		})
	}

	limitStorage := limitermocks.NewMockILimitStorage(t)
	limitStorage.EXPECT().GetLimitsByTypes(mock.AnythingOfType("[]string")).Return(&mockLimits, nil)

	return limitStorage
}
