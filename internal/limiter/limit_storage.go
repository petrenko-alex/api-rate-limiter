package limiter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrConnectFailed = errors.New("error connecting to db")

const DSN = "postgresql://main:main@localhost:5432/rate_limiter?sslmode=disable"

type LimitStorage struct {
	db  *sql.DB
	ctx context.Context
}

func NewLimitStorage() *LimitStorage {
	return &LimitStorage{}
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

func (s *LimitStorage) Connect(ctx context.Context) error {
	db, openErr := sql.Open("postgres", DSN)
	if openErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", openErr)
	}

	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		return fmt.Errorf(ErrConnectFailed.Error()+":%w", pingErr)
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *LimitStorage) Close(_ context.Context) error {
	closeErr := s.db.Close()
	if closeErr != nil {
		return closeErr
	}

	s.ctx = nil

	return nil
}

func (s *LimitStorage) scanRow(rows *sql.Rows) (*Limit, error) {
	limit := Limit{}
	nullableDescription := sql.NullString{}

	err := rows.Scan(&limit.limitType, &limit.value, &nullableDescription)
	if err != nil {
		return nil, err
	}

	if nullableDescription.Valid {
		limit.description = nullableDescription.String
	}

	return &limit, nil
}
