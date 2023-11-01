package limiter

type ILimitStorage interface {
	GetLimits() (*Limits, error)
}
