package limiter

import "github.com/petrenko-alex/api-rate-limiter/internal/ipnet"

const WhiteListLimiterIdentityKey = "ip"

// WhiteListLimiter лимитер на основе принципа белого списка.
// Если клиент есть в белом списке, то считается, что он всегда удовлетворяет лимитам.
type WhiteListLimiter struct {
	ruleService ipnet.IRuleService
}

func NewWhiteListLimiter(service *ipnet.RuleService) WhiteListLimiter {
	return WhiteListLimiter{
		ruleService: service,
	}
}

func (l WhiteListLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	ip, found := identity[WhiteListLimiterIdentityKey]
	if !found {
		return false, ErrIncorrectIdentity
	}

	inWhiteList, err := l.ruleService.InWhiteList(ip)
	if err != nil {
		return false, err
	}

	return inWhiteList, nil
}
