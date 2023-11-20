package limiter

import (
	"github.com/petrenko-alex/api-rate-limiter/internal/ipnet"
)

// LoginFormLimiter лимитер для использования в простых формах авторизации с логином и паролем.
type LoginFormLimiter struct {
	ruleService   ipnet.IRuleService
	bucketLimiter *CompositeBucketLimiter
}

func NewLoginFormLimiter(ruleService ipnet.IRuleService, limitStorage ILimitStorage, rate RefillRate) *LoginFormLimiter {
	return &LoginFormLimiter{
		ruleService:   ruleService,
		bucketLimiter: NewCompositeBucketLimiter(limitStorage, rate),
	}
}

func (l *LoginFormLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	validationErr := l.validateIdentity(identity)
	if validationErr != nil {
		return false, validationErr
	}

	inBlackList, blErr := l.ruleService.InBlackList(identity[IPLimit.String()])
	if inBlackList || blErr != nil {
		return false, blErr
	}

	inWhiteList, wlErr := l.ruleService.InWhiteList(identity[IPLimit.String()])
	if wlErr != nil {
		return false, wlErr
	}

	if inWhiteList {
		return true, nil
	}

	return l.bucketLimiter.SatisfyLimit(identity)
}

func (l *LoginFormLimiter) SetRequestCost(requestCost int) {
	l.bucketLimiter.SetRequestCost(requestCost)
}

func (l *LoginFormLimiter) validateIdentity(identity UserIdentityDto) error {
	if identity[IPLimit.String()] == "" ||
		identity[LoginLimit.String()] == "" ||
		identity[PasswordLimit.String()] == "" {
		// These identity keys are required
		return ErrIncorrectIdentity
	}

	return nil
}
