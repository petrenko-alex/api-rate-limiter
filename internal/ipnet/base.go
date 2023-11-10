package ipnet

type IRuleStorage interface {
	Create(rule Rule) (int, error)
	Delete(id int) error

	GetForIP(ip string) (*Rules, error)
	GetForType(ruleType RuleType) (*Rules, error)
}
