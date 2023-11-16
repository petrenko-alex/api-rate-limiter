package limiter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrConnectFailed = errors.New("error connecting to db")

type LimitStorage struct {
	db  *sql.DB
	ctx context.Context

	dsn string
}

func NewLimitStorage(dsn string) *LimitStorage {
	return &LimitStorage{dsn: dsn}
}

func (s *LimitStorage) GetLimits() (*Limits, error) {
	limits := make(Limits, 0)

	rows, err := s.db.QueryContext(s.ctx, "select type, value, description from rate_limit;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		limit, scanErr := s.scanRow(rows)
		if scanErr != nil {
			return nil, err
		}

		limits = append(limits, *limit)
	}

	return &limits, nil
}

func (s *LimitStorage) GetLimitsByTypes(types []string) (*Limits, error) {
	limits := make(Limits, 0)

	rows, err := s.db.QueryContext(
		s.ctx,
		"select type, value, description from rate_limit where type IN $1;",
		types,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		limit, scanErr := s.scanRow(rows)
		if scanErr != nil {
			return nil, err
		}

		limits = append(limits, *limit)
	}

	return &limits, nil
}

func (s *LimitStorage) Connect(ctx context.Context) error {
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

func (s *LimitStorage) Close(_ context.Context) error {
	if closeErr := s.db.Close(); closeErr != nil {
		return closeErr
	}

	s.ctx = nil

	return nil
}

func (s *LimitStorage) scanRow(rows *sql.Rows) (*Limit, error) {
	limit := Limit{}
	nullableDescription := sql.NullString{}

	err := rows.Scan(&limit.LimitType, &limit.Value, &nullableDescription)
	if err != nil {
		return nil, err
	}

	if nullableDescription.Valid {
		limit.Description = nullableDescription.String
	}

	return &limit, nil
}
