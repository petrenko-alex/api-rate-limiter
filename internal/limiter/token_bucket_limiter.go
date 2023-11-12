package limiter

const DefaultRequestCost = 1

// todo: doc and review

// TokenBucketLimiter лимитер с использованием алгоритма TokenBucket.
type TokenBucketLimiter struct {
	bucket      TokenBucket
	requestCost int
}

func NewTokenBucketLimiter(bucketKey string, bucketSize int, refillRate RefillRate) TokenBucketLimiter {
	return TokenBucketLimiter{
		bucket:      NewTokenBucket(bucketSize, refillRate),
		requestCost: DefaultRequestCost,
	}
}

func (l *TokenBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	l.bucket.Refill()

	if l.bucket.GetTokenCount() > 0 {
		l.bucket.GetToken(l.requestCost)

		return true, nil
	}

	return false, nil
}

func (l *TokenBucketLimiter) SetRequestCost(requestCost int) {
	l.requestCost = requestCost
}

func (l *TokenBucketLimiter) GetRequestsAllowed() int {
	l.bucket.Refill()

	return l.bucket.GetTokenCount() / l.requestCost
}
