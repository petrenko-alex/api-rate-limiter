package limiter

import (
	"errors"
	"time"
)

var (
	// ErrIncorrectIdentity Ошибка на случай некорректного входного аргумента identity.
	ErrIncorrectIdentity  = errors.New("not found appropriate key in user identity")
	ErrNotSupported       = errors.New("operation not supported")
	ErrIncorrectBucketKey = errors.New("incorrect bucket key")
)

// UserIdentityDto тип для идентификации клиента, запрос которого подвергается rate limit'ингу.
// Может содержать один или несколько пар ключ-значение. Лимитеры сами решают, с какими ключами работать.
type UserIdentityDto map[string]string

// ITokenBucket интерфейс хранилища токенов.
type ITokenBucket interface {
	GetSize() int
	GetTokenCount() int
	GetLastRefill() time.Time
	GetToken(int)
	Refill()
	Reset()
}

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

// ITokenBucketLimitService интерфейс лимитеров на основе TokenBucket.
type ITokenBucketLimitService interface {
	ILimitService

	SetRequestCost(int)
	GetRequestsAllowed(UserIdentityDto) (int, error)
	GetBuckets() map[string]*TokenBucket
	SweepBucket(string) error
}

// ITokenBucketGB сервис подчистки устаревших бакетов.
type ITokenBucketGB interface {
	Sweep()
}
