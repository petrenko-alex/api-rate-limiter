package limiter

import (
	"errors"
)

var ErrNoLimitsFound = errors.New("not found any limits for given identity")

// CompositeBucketLimiter лимитер с использованием нескольких bucket'ов
// Набор bucket'ов определяется на основе входных данных в UserIdentityDto (ключей).
// Объединение по логике И: для удовлетворения лимиту необходимо "пройти" все bucket'ы.
type CompositeBucketLimiter struct {
	limitStorage ILimitStorage

	limiters map[string]*TokenBucketLimiter

	refillRate  RefillRate
	requestCost int
}

func NewCompositeBucketLimiter(limitStorage ILimitStorage, refillRate RefillRate) *CompositeBucketLimiter {
	return &CompositeBucketLimiter{
		limitStorage: limitStorage,
		refillRate:   refillRate,
		requestCost:  DefaultRequestCost,
	}
}

func (l *CompositeBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	identityKeys := l.getIdentityKeys(identity)
	if len(identity) == 0 {
		return false, ErrIncorrectIdentity
	}

	if len(l.limiters) == 0 {
		limitersInitErr := l.initLimiters(identityKeys)
		if limitersInitErr != nil {
			return false, limitersInitErr
		}
	}

	for key := range identity {
		limiter, found := l.limiters[key]
		if !found {
			return false, ErrIncorrectIdentity
		}

		satisfies, checkErr := limiter.SatisfyLimit(identity)
		if checkErr != nil {
			return false, checkErr
		}

		if !satisfies {
			return false, nil // not satisfies if fails at least one limiter
		}
	}

	return true, nil // satisfies if pass all limiter
}

func (l *CompositeBucketLimiter) initLimiters(identityKeys []string) error {
	limits, getLimitsErr := l.limitStorage.GetLimitsByTypes(identityKeys)
	if getLimitsErr != nil || len(*limits) == 0 {
		return ErrNoLimitsFound
	}

	l.limiters = make(map[string]*TokenBucketLimiter, len(*limits))
	for _, limit := range *limits {
		limiterKey := limit.LimitType.String()
		limiter := NewTokenBucketLimiter(limiterKey, limit.Value, l.refillRate)
		limiter.SetRequestCost(l.requestCost)

		l.limiters[limiterKey] = limiter
	}

	return nil
}

func (l *CompositeBucketLimiter) getIdentityKeys(identity UserIdentityDto) []string {
	keys := make([]string, 0, len(identity))

	for key := range identity {
		keys = append(keys, key)
	}

	return keys
}

func (l *CompositeBucketLimiter) SetRequestCost(requestCost int) {
	l.requestCost = requestCost

	if l.limiters == nil || len(l.limiters) == 0 {
		return
	}

	for _, limiter := range l.limiters {
		limiter.SetRequestCost(requestCost)
	}
}
