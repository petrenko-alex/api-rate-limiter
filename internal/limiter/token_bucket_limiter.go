package limiter

import "errors"

const DefaultRequestCost = 1

var ErrIncorrectIdentity = errors.New("not found bucket key in user identity")

// todo: doc and review

// TokenBucketLimiter лимитер с использованием алгоритма TokenBucket.
type TokenBucketLimiter struct {
	buckets          map[string]*TokenBucket
	bucketSize       int
	bucketRefillRate RefillRate

	requestCost int
	bucketKey   string
}

func NewTokenBucketLimiter(bucketKey string, bucketSize int, refillRate RefillRate) TokenBucketLimiter {
	return TokenBucketLimiter{
		buckets:     make(map[string]*TokenBucket),
		requestCost: DefaultRequestCost,

		bucketKey:        bucketKey,
		bucketSize:       bucketSize,
		bucketRefillRate: refillRate,
	}
}

func (l *TokenBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	identityValue, found := identity[l.bucketKey]
	if !found {
		return false, ErrIncorrectIdentity
	}

	bucket := l.initBucket(identityValue)

	bucket.Refill()

	if bucket.GetTokenCount() > 0 {
		bucket.GetToken(l.requestCost)

		return true, nil
	}

	return false, nil
}

func (l *TokenBucketLimiter) SetRequestCost(requestCost int) {
	l.requestCost = requestCost
}

func (l *TokenBucketLimiter) GetRequestsAllowed(identity UserIdentityDto) (int, error) {
	identityValue, found := identity[l.bucketKey]
	if !found {
		return 0, ErrIncorrectIdentity
	}

	bucket := l.initBucket(identityValue)

	bucket.Refill()

	return bucket.GetTokenCount() / l.requestCost, nil
}

func (l *TokenBucketLimiter) initBucket(identityValue string) *TokenBucket {
	bucket := l.findBucket(identityValue)
	if bucket == nil {
		bucket = l.createBucket(identityValue)
	}

	return bucket
}

func (l *TokenBucketLimiter) findBucket(identityValue string) *TokenBucket {
	bucket, found := l.buckets[identityValue]
	if !found {
		return nil
	}

	return bucket
}

func (l *TokenBucketLimiter) createBucket(identityValue string) *TokenBucket {
	newBucket := NewTokenBucket(l.bucketSize, l.bucketRefillRate)
	l.buckets[identityValue] = &newBucket

	return l.buckets[identityValue]
}
