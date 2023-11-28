package main

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up2, Down2)
}

func Up2(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		ALTER TABLE events ADD COLUMN notified boolean;
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down2(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "ALTER TABLE events DROP COLUMN notified")
	if err != nil {
		return err
	}
	return nil
}
