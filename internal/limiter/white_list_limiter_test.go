package limiter_test //nolint:dupl

import (
	"testing"

	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
	ipnetmocks "github.com/petrenko-alex/api-rate-limiter/internal/ipnet/mocks"
	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"github.com/stretchr/testify/require"
)

func TestWhiteListLimiter_SatisfyLimit(t *testing.T) {
	rules := ipnet.Rules{
		ipnet.Rule{
			ID:       1,
			IP:       "192.168.1.1",
			RuleType: ipnet.WhiteList,
		},
		ipnet.Rule{
			ID:       2,
			IP:       "125.130.2.3",
			RuleType: ipnet.WhiteList,
		},
		ipnet.Rule{
			ID:       3,
			IP:       "192.168.3.0/24",
			RuleType: ipnet.WhiteList,
		},
	}
	mockStorage := ipnetmocks.NewMockIRuleStorage(t)
	mockStorage.EXPECT().GetForType(ipnet.WhiteList).Return(&rules, nil)

	whiteListLimiter := limiter.NewWhiteListLimiter(
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
			expected: true,
		},
		{
			name:     "found by subnet match",
			ip:       "192.168.3.25",
			expected: true,
		},
		{
			name:     "not found",
			ip:       "222.113.10.3",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			satisfies, err := whiteListLimiter.SatisfyLimit(limiter.UserIdentityDto{"ip": test.ip})

			require.NoError(t, err)
			require.Equal(t, test.expected, satisfies)
		})
	}
}

func TestWhiteListLimiter_SatisfyLimit_Error(t *testing.T) {
	t.Run("incorrect identity error", func(t *testing.T) {
		mockStorage := ipnetmocks.NewMockIRuleStorage(t)
		whiteListLimiter := limiter.NewWhiteListLimiter(
			ipnet.NewRuleService(mockStorage),
		)

		identity := limiter.UserIdentityDto{"login": "admin"} // white list limiter needs ip
		_, err := whiteListLimiter.SatisfyLimit(identity)
		require.ErrorIs(t, err, limiter.ErrIncorrectIdentity)
	})
}
