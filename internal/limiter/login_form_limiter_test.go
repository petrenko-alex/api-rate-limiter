package limiter_test

import (
	"testing"

	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestLoginFormLimiter_SatisfyLimit(t *testing.T) {
	identity := limiter.UserIdentityDto{
		limiter.IPLimit.String():       "192.168.1.1",
		limiter.LoginLimit.String():    "lucky",
		limiter.PasswordLimit.String(): "root",
	}

	t.Run("ip in white list", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter()

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// loginFormLimiter.SetRequestCost(4)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("ip in black list", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter()
		// loginFormLimiter.SetRequestCost(4)

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		// loginFormLimiter.SetRequestCost(4)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

	})

	t.Run("ip in both lists", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter()
		identity[limiter.IPLimit.String()] = "10.9.123.12"

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("no ip in lists", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter()
		identity[limiter.IPLimit.String()] = "unknown"

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		// loginFormLimiter.SetRequestCost(4)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})
}

func TestLoginFormLimiter_SatisfyLimit_Error(t *testing.T) {
	t.Run("incorrect identity", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter()

		_, err := loginFormLimiter.SatisfyLimit(limiter.UserIdentityDto{
			limiter.LoginLimit.String():    "lucky",
			limiter.PasswordLimit.String(): "root",
		})
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)

		_, err = loginFormLimiter.SatisfyLimit(limiter.UserIdentityDto{
			limiter.PasswordLimit.String(): "root",
		})
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)

		_, err = loginFormLimiter.SatisfyLimit(limiter.UserIdentityDto{})
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}
