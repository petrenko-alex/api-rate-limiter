package ipnet

type RuleType string

const (
	BlackList RuleType = "black"
	WhiteList RuleType = "white"
)

type Rules []Rule

type Rule struct {
	ID       int
	IP       string
	RuleType RuleType
}
