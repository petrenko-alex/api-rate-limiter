package limiter

// BlackListLimiter лимитер на основе принципа черного списка.
// Если клиент есть в черном списке, то считается, что он всегда не попадает в лимит.
type BlackListLimiter struct {
}

func (l BlackListLimiter) SatisfyLimit(identity UserIdentityDto) (bool, error) {
	// Get ip from identity
	// 		error if no ip
	// Use RuleService::InBlackList() to check.
	// If TRUE - return FALSE
	// IF FALSE - return TRUE

	panic("implement me")
}
