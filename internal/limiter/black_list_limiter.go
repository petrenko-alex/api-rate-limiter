package limiter

import (
	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
)

const BlackListLimiterIdentityKey = "ip"

// BlackListLimiter лимитер на основе принципа черного списка.
// Если клиент есть в черном списке, то считается, что он всегда не попадает в лимит.
type BlackListLimiter struct {
	ruleService ipnet.IRuleService
}

func NewBlackListLimiter(service *ipnet.RuleService) *BlackListLimiter {
	return &BlackListLimiter{
		ruleService: service,
	}
}

func (l BlackListLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	ip, found := identity[BlackListLimiterIdentityKey]
	if !found {
		return false, ErrIncorrectIdentity
	}

	inBlackList, err := l.ruleService.InBlackList(ip)
	if err != nil {
		return false, err
	}

	return !inBlackList, nil
}

func (l BlackListLimiter) ResetLimit(_ UserIdentityDto) error {
	return ErrNotSupported
}
