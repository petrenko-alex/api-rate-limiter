package limiter

// TokenBucketLimiter лимитер с использованием алгоритма TokenBucket.
type TokenBucketLimiter struct {
}

func NewTokenBucketLimiter() TokenBucketLimiter {
	return TokenBucketLimiter{}
}

func (l TokenBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	// Realize TokenBucket algo

	panic("implement me")
}
