package ipnet

import (
	"errors"
	"net"
)

var (
	ErrRuleNotFound   = errors.New("rule not found")
	ErrInvalidInputIP = errors.New("incorrect IP passed")
)

type RuleService struct {
	ruleStorage IRuleStorage
}

func NewRuleService(ruleStorage IRuleStorage) *RuleService {
	return &RuleService{ruleStorage: ruleStorage}
}

func (s RuleService) InWhiteList(ip string) (bool, error) {
	return s.inList(ip, WhiteList)
}

func (s RuleService) InBlackList(ip string) (bool, error) {
	return s.inList(ip, BlackList)
}

func (s RuleService) WhiteListAdd(ip string) error {
	return s.listAdd(ip, WhiteList)
}

func (s RuleService) WhiteListDelete(ip string) error {
	return s.listDelete(ip, WhiteList)
}

func (s RuleService) BlackListAdd(ip string) error {
	return s.listAdd(ip, BlackList)
}

func (s RuleService) BlackListDelete(ip string) error {
	return s.listDelete(ip, BlackList)
}

func (s RuleService) listAdd(ip string, listType RuleType) error {
	_, err := s.ruleStorage.Create(Rule{
		IP:       ip,
		RuleType: listType,
	})

	return err
}

func (s RuleService) listDelete(ip string, listType RuleType) error {
	rules, err := s.ruleStorage.Find(ip, listType)
	if err != nil {
		return err
	}

	if len(*rules) == 0 {
		return ErrRuleNotFound
	}

	for _, rule := range *rules {
		deleteErr := s.ruleStorage.Delete(rule.ID)
		if deleteErr != nil {
			return deleteErr
		}
	}

	return nil
}

func (s RuleService) inList(ip string, listType RuleType) (bool, error) {
	rules, getErr := s.ruleStorage.GetForType(listType)
	if getErr != nil {
		return false, getErr
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false, ErrInvalidInputIP
	}

	if len(*rules) == 0 {
		return false, nil
	}

	for _, rule := range *rules {
		ruleParsedIP := net.ParseIP(rule.IP)
		if ruleParsedIP != nil && parsedIP.Equal(ruleParsedIP) { // direct ip match
			return true, nil
		}

		_, ipNet, parseErr := net.ParseCIDR(rule.IP)
		if parseErr != nil {
			continue
		}

		if ipNet.Contains(parsedIP) { // subnet match
			return true, nil
		}
	}

	return false, nil
}
