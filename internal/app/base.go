package app

type IApp interface {
	LimitCheck(ip, login, password string) (bool, error)

	WhiteListAdd(ip string) error
	WhiteListDelete(ip string) error

	BlackListAdd(ip string) error
	BlackListDelete(ip string) error
}
