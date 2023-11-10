package limiter

// WhiteListLimiter лимитер на основе принципа белого списка.
// Если клиент есть в белом списке, то считается, что он всегда удовлетворяет лимитам.
type WhiteListLimiter struct {
}

func (l WhiteListLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	// Get ip from identity
	// 		error if no ip
	// Use RuleService::InWhiteList() to check.
	// If TRUE - return TRUE
	// IF FALSE - return FALSE

	panic("implement me")
}
