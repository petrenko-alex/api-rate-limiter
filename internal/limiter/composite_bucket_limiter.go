package limiter

import "errors"

// todo: rename to CompositeTokenBucketLimiter ?

var ErrNoRulesFound = errors.New("not found any rules for given identity")

// CompositeBucketLimiter лимитер с использованием нескольких bucket'ов
// Набор bucket'ов определяется на основе входных данных в UserIdentityDto (ключей).
// Объединение по логике И: для удовлетворения лимиту необходимо "пройти" все bucket'ы.
type CompositeBucketLimiter struct{}

func NewCompositeBucketLimiter() CompositeBucketLimiter {
	return CompositeBucketLimiter{}
}

func (l CompositeBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	// Get ip from identity ??
	// 		error if no ip
	// Get all keys from identity
	// Find all rules for keys and get limits
	// 		What if not all keys were found? (fail open vs fail close)
	// Loop for identity keys
	// 		Find bucket by key
	//		IF NOT FOUND - create bucket
	//		IF FOUND - check bucket
	//			IF SATISFY - go to identity key/bucket
	// 			IF NOT SATISFY - return FALSE
	// return true (because all buckets satisfy)

	return false, nil
}
