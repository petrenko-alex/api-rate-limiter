package limiter

// LoginFormLimiter лимитер для использования в простых формах авторизации с логином и паролем.
type LoginFormLimiter struct {
}

func (l LoginFormLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	// Get ip, login and password from identity
	// 		error if no all data found
	// Check ip in whitelist
	// 		IF FOUND - return TRUE
	// Check ip in blacklist
	//		IF FOUND - return FALSE
	// Call CompositeBucketLimiter to check limits using TokenBucket's

	panic("implement me")
}
