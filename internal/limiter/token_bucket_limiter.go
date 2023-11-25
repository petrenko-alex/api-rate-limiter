package limiter

const DefaultRequestCost = 1

// TokenBucketLimiter позволяет задать rate limit для запросов с использованием алгоритма TokenBucket.
//
// Для идентификации клиента запроса используется обобщенный объект UserIdentityDto.
// Ключ bucketKey используется для поиска идентификатора клиента в UserIdentityDto.
//
// Позволяет проверять возможность выполнения очередного запроса и получать количество доступных.
type TokenBucketLimiter struct {
	buckets    map[string]*TokenBucket
	bucketSize int

	// Скорость пополнения токенов корзины.
	bucketRefillRate RefillRate

	// Количество токенов, которое тратится на один запрос при вызове SatisfyLimit.
	// По умолчанию - 1.
	requestCost int

	// Ключ корзины. По данному ключу происходит поиск идентификационных данных методом SatisfyLimit и ResetLimit.
	bucketKey string
}

func NewTokenBucketLimiter(bucketKey string, bucketSize int, refillRate RefillRate) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:     make(map[string]*TokenBucket),
		requestCost: DefaultRequestCost,

		bucketKey:        bucketKey,
		bucketSize:       bucketSize,
		bucketRefillRate: refillRate,
	}
}

// SatisfyLimit проверяет возможность выполнения запроса для identity c учетом текущей стоимости запроса.
// Происходит забор токенов из корзины для identity c учетом текущей стоимости запроса.
// Выполняется пополнение корзины для identity с учетом заданной скорости пополнения.
func (l *TokenBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	identityValue, found := identity[l.bucketKey]
	if !found {
		return false, ErrIncorrectIdentity
	}

	bucket := l.initBucket(identityValue)

	bucket.Refill()

	if bucket.GetTokenCount() > 0 && l.requestCost <= bucket.GetTokenCount() {
		bucket.GetToken(l.requestCost)

		return true, nil
	}

	return false, nil
}

func (l *TokenBucketLimiter) ResetLimit(identity UserIdentityDto) error {
	identityValue, foundIdentityValue := identity[l.bucketKey]
	if !foundIdentityValue {
		return ErrIncorrectIdentity
	}

	bucket, foundBucket := l.buckets[identityValue]
	if !foundBucket {
		return nil
	}

	bucket.Reset()

	return nil
}

func (l *TokenBucketLimiter) SweepBucket(bucketKey string) error {
	delete(l.buckets, bucketKey)

	return nil
}

func (l *TokenBucketLimiter) SetRequestCost(requestCost int) {
	l.requestCost = requestCost
}

// GetRequestsAllowed возвращает количество возможных запросов для identity с учетом текущей стоимости запроса.
func (l *TokenBucketLimiter) GetRequestsAllowed(identity UserIdentityDto) (int, error) {
	identityValue, found := identity[l.bucketKey]
	if !found {
		return 0, ErrIncorrectIdentity
	}

	bucket := l.initBucket(identityValue)

	bucket.Refill()

	return bucket.GetTokenCount() / l.requestCost, nil
}

func (l *TokenBucketLimiter) GetBuckets() map[string]*TokenBucket {
	return l.buckets
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
