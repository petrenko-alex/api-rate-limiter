package limiter_test

import (
	"testing"

	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	ipnetmocks "github.com/petrenko-alex/api-rate-limiter/internal/ipnet/mocks"
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestBlackListLimiter_SatisfyLimit(t *testing.T) {
	rules := ipnet.Rules{
		ipnet.Rule{
			ID:       1,
			IP:       "192.168.1.1",
			RuleType: ipnet.BlackList,
		},
		ipnet.Rule{
			ID:       2,
			IP:       "125.130.2.3",
			RuleType: ipnet.BlackList,
		},
		ipnet.Rule{
			ID:       3,
			IP:       "192.168.3.0/24",
			RuleType: ipnet.BlackList,
		},
	}
	mockStorage := ipnetmocks.NewMockIRuleStorage(t)
	mockStorage.EXPECT().GetForType(ipnet.BlackList).Return(&rules, nil)

	blackListLimiter := limiter.NewBlackListLimiter(
		ipnet.NewRuleService(mockStorage),
	)

	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "found by exact ip match",
			ip:       "192.168.1.1",
			expected: false,
		},
		{
			name:     "found by subnet match",
			ip:       "192.168.3.25",
			expected: false,
		},
		{
			name:     "not found",
			ip:       "222.113.10.3",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			satisfies, err := blackListLimiter.SatisfyLimit(limiter.UserIdentityDto{"ip": test.ip})

			require.NoError(t, err)
			require.Equal(t, test.expected, satisfies)
		})
	}
}

func TestBlackListLimiter_SatisfyLimit_Error(t *testing.T) {
	t.Run("incorrect identity error", func(t *testing.T) {
		mockStorage := ipnetmocks.NewMockIRuleStorage(t)
		blackListLimiter := limiter.NewBlackListLimiter(
			ipnet.NewRuleService(mockStorage),
		)

		identity := limiter.UserIdentityDto{"login": "admin"} // black list limiter needs ip
		_, err := blackListLimiter.SatisfyLimit(identity)
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}
