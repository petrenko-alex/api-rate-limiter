package limiter_test

import (
	"testing"

	"api-rate-limiter/internal/ipnet"
	"api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

// todo: use table test
func TestWhiteListLimiter_SatisfyLimit(t *testing.T) {
	strg := ipnet.NewRuleStorage("dsn") // todo: mock
	srvc := ipnet.NewRuleService(strg)

	whiteListLimiter := limiter.NewWhiteListLimiter(srvc)

	t.Run("found by exact ip match", func(t *testing.T) {
		satisfies, err := whiteListLimiter.SatisfyLimit(limiter.UserIdentityDto{"ip": "192.168.1.1"})

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("found by subnet match", func(t *testing.T) {
		satisfies, err := whiteListLimiter.SatisfyLimit(limiter.UserIdentityDto{"ip": "125.110.13.25"})

		require.NoError(t, err)
		require.True(t, satisfies)
	})

	t.Run("not found", func(t *testing.T) {
		satisfies, err := whiteListLimiter.SatisfyLimit(limiter.UserIdentityDto{"ip": "222.113.10.3"})

		require.NoError(t, err)
		require.False(t, satisfies)
	})
}

func TestWhiteListLimiter_SatisfyLimit_Error(t *testing.T) {
	t.Run("incorrect identity error", func(t *testing.T) {
		whiteListLimiter := limiter.NewWhiteListLimiter()

		identity := limiter.UserIdentityDto{"login": "admin"} // white list limiter needs ip
		_, err := whiteListLimiter.SatisfyLimit(identity)
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}
