package limiter

import "errors"

// ErrIncorrectIdentity Ошибка на случай некорректного входного аргумента identity
var ErrIncorrectIdentity = errors.New("not found appropriate key in user identity")

type UserIdentityDto map[string]string

type ILimitStorage interface {
	GetLimits() (*Limits, error)
}

type ILimitService interface {
	SatisfyLimit(UserIdentityDto) (bool, error)
}
