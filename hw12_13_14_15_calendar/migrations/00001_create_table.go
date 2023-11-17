package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(Up, Down)
}

func Up(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE events
		(
			id            bpchar      NOT NULL,
			title         text        NOT NULL,
			start_dt      timestamp   NOT NULL,
			end_dt        timestamp   NOT NULL,
			description   text,
			user_id       bpchar      NOT NULL,
			notify_before interval    second (0),
			CONSTRAINT uqnique_id UNIQUE (id),
			CONSTRAINT pk_event_id PRIMARY KEY (id)
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func Down(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE events")
	if err != nil {
		return err
	}
	return nil
}
