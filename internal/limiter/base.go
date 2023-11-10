package limiter

type UserIdentityDto map[string]string

type ILimitStorage interface {
	GetLimits() (*Limits, error)
}

type ILimitService interface {
	SatisfyLimit(UserIdentityDto) (bool, error)
}
