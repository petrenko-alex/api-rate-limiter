package migrations

import (
	"context"
	"database/sql"

	"github.com/petrenko-alex/api-rate-limiter/internal/limiter"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upFillLimitTable, downFillLimitTable)
}

func upFillLimitTable(ctx context.Context, tx *sql.Tx) error {
	query := `insert into rate_limit(type, value, description) 
				values ($1, 10, 'Ограничение для логина'),
				       ($2, 100, 'Ограничение для пароля'),
				       ($3, 1000, null)
	;`

	if _, err := tx.ExecContext(ctx, query, limiter.LoginLimit, limiter.PasswordLimit, limiter.IPLimit); err != nil {
		return err
	}

	return nil
}

func downFillLimitTable(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `truncate table rate_limit;`); err != nil {
		return err
	}

	return nil
}
