package limiter

type LimitType string

const (
	LoginLimit    LimitType = "login"
	PasswordLimit LimitType = "password"
	IPLimit       LimitType = "ip"
)

type Limits []Limit

type Limit struct {
	limitType   LimitType
	value       string
	description string
}
