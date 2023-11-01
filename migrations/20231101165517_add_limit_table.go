package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLimitTable, downAddLimitTable)
}

func upAddLimitTable(ctx context.Context, tx *sql.Tx) error {
	query := `create table rate_limit(
    type varchar(50) primary key, 
    value int not null,
    description varchar(255) null 
);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddLimitTable(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS rate_limit;"); err != nil {
		return err
	}

	return nil
}
