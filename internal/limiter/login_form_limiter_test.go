package limiter_test

import (
	"testing"
	"time"

	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet/mocks"
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	limitermocks "github.com/petrenko-alex/api-rate-limiter/internal/limiter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoginFormLimiter_SatisfyLimit(t *testing.T) {
	limit := 3
	refillRate := limiter.NewRefillRate(limit, time.Second*1)
	whiteListIP, blackListIP, bothListIP, unknownIP := "192.168.1.1", "192.150.10.3", "10.9.123.12", "5.5.5.5"
	identity := limiter.UserIdentityDto{
		limiter.IPLimit.String():       "unknown",
		limiter.LoginLimit.String():    "lucky",
		limiter.PasswordLimit.String(): "root",
	}

	// Mock RuleService
	ruleStorage := mocks.NewMockIRuleStorage(t)
	ruleStorage.EXPECT().GetForType(ipnet.WhiteList).Return(&ipnet.Rules{
		ipnet.Rule{ID: 1, IP: whiteListIP, RuleType: ipnet.WhiteList},
		ipnet.Rule{ID: 3, IP: bothListIP, RuleType: ipnet.WhiteList},
	}, nil).Maybe()
	ruleStorage.EXPECT().GetForType(ipnet.BlackList).Return(&ipnet.Rules{
		ipnet.Rule{ID: 2, IP: blackListIP, RuleType: ipnet.BlackList},
		ipnet.Rule{ID: 4, IP: bothListIP, RuleType: ipnet.BlackList},
	}, nil).Maybe()
	ruleService := ipnet.NewRuleService(ruleStorage)

	// Mock LimitStorage
	limitStorage := limitermocks.NewMockILimitStorage(t)
	limitStorage.EXPECT().GetLimitsByTypes(mock.AnythingOfType("[]string")).Return(&limiter.Limits{
		limiter.Limit{LimitType: limiter.IPLimit, Value: limit},
		limiter.Limit{LimitType: limiter.LoginLimit, Value: limit},
		limiter.Limit{LimitType: limiter.PasswordLimit, Value: limit},
	}, nil).Maybe()

	t.Run("ip in white list", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter(
			ruleService,
			limiter.NewCompositeBucketLimiter(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = whiteListIP

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		loginFormLimiter.SetRequestCost(limit + 1)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("ip in black list", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter(
			ruleService,
			limiter.NewCompositeBucketLimiter(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = blackListIP
		loginFormLimiter.SetRequestCost(limit + 1)

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)

		loginFormLimiter.SetRequestCost(1)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("ip in both lists", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter(
			ruleService,
			limiter.NewCompositeBucketLimiter(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = bothListIP

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})

	t.Run("no ip in lists", func(t *testing.T) {
		loginFormLimiter := limiter.NewLoginFormLimiter(
			ruleService,
			limiter.NewCompositeBucketLimiter(limitStorage, refillRate),
		)
		identity[limiter.IPLimit.String()] = unknownIP

		satisfies, err := loginFormLimiter.SatisfyLimit(identity)
		require.True(t, satisfies)
		require.NoError(t, err)

		loginFormLimiter.SetRequestCost(limit + 1)
		satisfies, err = loginFormLimiter.SatisfyLimit(identity)
		require.False(t, satisfies)
		require.NoError(t, err)
	})
}

func TestLoginFormLimiter_SatisfyLimit_Error(t *testing.T) {
	ruleService := ipnet.NewRuleService(mocks.NewMockIRuleStorage(t))
	limitStorage := limitermocks.NewMockILimitStorage(t)
	refillRate := limiter.NewRefillRate(3, time.Second*1)

	t.Run("incorrect identity", func(t *testing.T) {
		expectedErr := limiter.ErrIncorrectIdentity
		emptyIdentity := limiter.UserIdentityDto{}
		notFullIdentity := limiter.UserIdentityDto{
			limiter.LoginLimit.String():    "lucky",
			limiter.PasswordLimit.String(): "root",
		}
		loginFormLimiter := limiter.NewLoginFormLimiter(
			ruleService,
			limiter.NewCompositeBucketLimiter(limitStorage, refillRate),
		)

		// not full identity #1
		_, err := loginFormLimiter.SatisfyLimit(notFullIdentity)
		require.ErrorIs(t, err, expectedErr)

		// not full identity #2
		delete(notFullIdentity, limiter.LoginLimit.String())
		_, err = loginFormLimiter.SatisfyLimit(notFullIdentity)
		require.ErrorIs(t, err, expectedErr)

		// empty identity
		_, err = loginFormLimiter.SatisfyLimit(emptyIdentity)
		require.ErrorIs(t, err, expectedErr)
	})
}
