package postgres

import (
	"context"
	"database/sql"
	_ "embed"
)

//go:embed migrations/001_init.sql
var initSQL string

func Migrate(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, initSQL)
	return err
}
