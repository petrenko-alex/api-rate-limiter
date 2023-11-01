package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddIpNetRuleTable, downAddIpNetRuleTable)
}

func upAddIpNetRuleTable(ctx context.Context, tx *sql.Tx) error {
	query := `CREATE TABLE ip_net_rule(
    id int not null primary key,
    ip varchar(45) not null,
    net int not null,
    type varchar(255) not null
);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downAddIpNetRuleTable(ctx context.Context, tx *sql.Tx) error {
	query := `DROP TABLE IF EXISTS ip_net_rule`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}
