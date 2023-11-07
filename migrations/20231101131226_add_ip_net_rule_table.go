package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddIPNetRuleTable, downAddIPNetRuleTable)
}

func upAddIPNetRuleTable(ctx context.Context, tx *sql.Tx) error {
	query := `create table ip_net_rule(
    id bigint generated always as identity primary key,
    ip varchar(50) not null,
    type varchar(255) not null
);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddIPNetRuleTable(ctx context.Context, tx *sql.Tx) error {
	query := `drop table if exists ip_net_rule`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
