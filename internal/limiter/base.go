package limiter

import "errors"

var (
	// ErrIncorrectIdentity Ошибка на случай некорректного входного аргумента identity.
	ErrIncorrectIdentity = errors.New("not found appropriate key in user identity")
	ErrNotSupported      = errors.New("operation not supported")
)

// UserIdentityDto тип для идентификации клиента, запрос которого подвергается rate limit'ингу.
// Может содержать один или несколько пар ключ-значение. Лимитеры сами решают, с какими ключами работать.
type UserIdentityDto map[string]string

// ILimitStorage хранилище лимитов (правил) rate limit'инга запросов.
type ILimitStorage interface {
	GetLimits() (*Limits, error)
	GetLimitsByTypes([]string) (*Limits, error)
}

// ILimitService основной сервис проверки запроса на rate limit.
type ILimitService interface {
	SatisfyLimit(UserIdentityDto) (bool, error)
	ResetLimit(UserIdentityDto) error
}

// ITokenBucketGB сервис подчистки устаревших бакетов.
type ITokenBucketGB interface {
	Sweep()
}
