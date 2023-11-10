package limiter

// CompositeBucketLimiter лимитер с использованием нескольких bucket'ов
// Набор bucket'ов определяется на основе входных данных в UserIdentityDto.
// Объединение по логике И: для удовлетворения лимиту необходимо "пройти" все bucket'ы.
type CompositeBucketLimiter struct {
}

func (l CompositeBucketLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	// Get ip from identity
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

	panic("implement me")
}
