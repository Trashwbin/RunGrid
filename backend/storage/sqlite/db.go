package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS groups (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	display_order INTEGER NOT NULL DEFAULT 0,
	color TEXT NOT NULL DEFAULT '',
	category TEXT NOT NULL DEFAULT 'app',
	icon TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS items (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	path TEXT NOT NULL,
	target_name TEXT NOT NULL DEFAULT '',
	type TEXT NOT NULL,
	icon_path TEXT NOT NULL DEFAULT '',
	group_id TEXT NOT NULL DEFAULT '',
	tags TEXT NOT NULL DEFAULT '[]',
	favorite INTEGER NOT NULL DEFAULT 0,
	launch_count INTEGER NOT NULL DEFAULT 0,
	last_used_at INTEGER,
	hidden INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_items_group ON items(group_id);
CREATE INDEX IF NOT EXISTS idx_items_name ON items(name);
CREATE INDEX IF NOT EXISTS idx_items_path ON items(path);
CREATE INDEX IF NOT EXISTS idx_items_last_used ON items(last_used_at);
`

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, nil
}

func EnsureSchema(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return err
	}

	if err := ensureItemColumns(ctx, db); err != nil {
		return err
	}
	return ensureGroupColumns(ctx, db)
}

func ensureItemColumns(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "PRAGMA table_info(items)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasTargetName := false
	for rows.Next() {
		var (
			cid       int
			name      string
			colType   string
			notNull   int
			dfltValue sql.NullString
			pk        int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "target_name" {
			hasTargetName = true
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if !hasTargetName {
		if _, err := db.ExecContext(ctx, "ALTER TABLE items ADD COLUMN target_name TEXT NOT NULL DEFAULT ''"); err != nil {
			return err
		}
	}

	if _, err := db.ExecContext(ctx, "CREATE INDEX IF NOT EXISTS idx_items_target_name ON items(target_name)"); err != nil {
		return err
	}

	return nil
}

func ensureGroupColumns(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "PRAGMA table_info(groups)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasCategory := false
	hasIcon := false
	for rows.Next() {
		var (
			cid       int
			name      string
			colType   string
			notNull   int
			dfltValue sql.NullString
			pk        int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "category" {
			hasCategory = true
		}
		if name == "icon" {
			hasIcon = true
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	if !hasCategory {
		if _, err := db.ExecContext(ctx, "ALTER TABLE groups ADD COLUMN category TEXT NOT NULL DEFAULT 'app'"); err != nil {
			return err
		}
	}

	if !hasIcon {
		if _, err := db.ExecContext(ctx, "ALTER TABLE groups ADD COLUMN icon TEXT NOT NULL DEFAULT ''"); err != nil {
			return err
		}
	}

	return nil
}
