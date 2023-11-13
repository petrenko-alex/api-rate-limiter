package ipnet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq" // driver import
)

var ErrConnectFailed = errors.New("error connecting to db")

type RuleStorage struct {
	db  *sql.DB
	ctx context.Context

	dsn string
}

func NewRuleStorage(dsn string) *RuleStorage {
	return &RuleStorage{dsn: dsn}
}

func (s *RuleStorage) Create(rule Rule) (int, error) {
	err := s.db.QueryRowContext(
		s.ctx,
		"INSERT INTO ip_net_rule(ip, type) VALUES ($1, $2) RETURNING id;",
		rule.IP,
		rule.RuleType,
	).Scan(&rule.ID)
	if err != nil {
		return 0, err
	}

	return rule.ID, nil
}

func (s *RuleStorage) Delete(id int) error {
	_, err := s.db.ExecContext(s.ctx, "DELETE FROM ip_net_rule WHERE id=$1;", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *RuleStorage) GetForIP(ip string) (*Rules, error) {
	rules := Rules{}

	rows, err := s.db.QueryContext(
		s.ctx,
		"SELECT id, ip, type FROM ip_net_rule WHERE ip=$1;",
		ip,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rule := Rule{}

		err = rows.Scan(&rule.ID, &rule.IP, &rule.RuleType)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return &rules, nil
}

func (s *RuleStorage) GetForType(ruleType RuleType) (*Rules, error) {
	rules := make(Rules, 0)

	rows, err := s.db.QueryContext(
		s.ctx,
		"SELECT id, ip, type FROM ip_net_rule WHERE type=$1;",
		ruleType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rule := Rule{}
		err = rows.Scan(&rule.ID, &rule.IP, &rule.RuleType)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return &rules, nil
}

func (s *RuleStorage) Find(ip string, ruleType RuleType) (*Rules, error) {
	rules := make(Rules, 0)

	rows, err := s.db.QueryContext(
		s.ctx,
		"SELECT id, ip, type FROM ip_net_rule WHERE ip=$1 AND type=$2;",
		ip,
		ruleType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rule := Rule{}

		scanErr := rows.Scan(&rule.ID, &rule.IP, &rule.RuleType)
		if scanErr != nil {
			return nil, scanErr
		}

		rules = append(rules, rule)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return &rules, nil
}

func (s *RuleStorage) Connect(ctx context.Context) error {
	db, openErr := sql.Open("postgres", s.dsn)
	if openErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", openErr)
	}

	if pingErr := db.PingContext(ctx); pingErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", pingErr)
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *RuleStorage) Close(_ context.Context) error {
	if closeErr := s.db.Close(); closeErr != nil {
		return closeErr
	}

	s.ctx = nil

	return nil
}
