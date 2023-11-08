package ipnet

type IRuleStorage interface {
	Create(rule Rule) (int, error)
	Delete(id int) error

	GetForIP(ip string) (*Rules, error)
	GetForType(ruleType RuleType) (*Rules, error)
	Find(ip string, ruleType RuleType) (*Rules, error)
}

type IRuleService interface {
	InWhiteList(ip string) (bool, error)
	InBlackList(ip string) (bool, error)

	WhiteListAdd(ip string) error
	WhiteListDelete(ip string) error

	BlackListAdd(ip string) error
	BlackListDelete(ip string) error
}
