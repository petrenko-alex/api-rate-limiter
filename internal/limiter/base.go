package limiter

import "errors"

// ErrIncorrectIdentity Ошибка на случай некорректного входного аргумента identity.
var ErrIncorrectIdentity = errors.New("not found appropriate key in user identity")

// UserIdentityDto тип для идентификации клиента, запрос которого подвергается rate limit'ингу.
// Может содержать один или несколько пар ключ-значение. Лимитеры сами решают, с какими ключами работать.
type UserIdentityDto map[string]string

// ILimitStorage хранилище лимитов (правил) rate limit'инга запросов.
type ILimitStorage interface {
	GetLimits() (*Limits, error)
}

// ILimitService основной сервис проверки запроса на rate limit.
type ILimitService interface {
	SatisfyLimit(UserIdentityDto) (bool, error)
}
