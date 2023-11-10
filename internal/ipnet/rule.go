package ipnet

type RuleType uint8

const (
	BlackList RuleType = iota
	WhiteList
)

type Rules []Rule

type Rule struct {
	ID       int
	IP       string
	RuleType RuleType
}
