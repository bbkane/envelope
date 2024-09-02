package sqlite

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"time"

	_ "modernc.org/sqlite"
)

// NOTE: most of this is heavily inspired by https://github.com/benbjohnson/wtf/blob/05bc90c940d5f9e2490fc93cf467d9e8aa48ad63/sqlite/sqlite.go

//go:embed migrations/*.sql
var migrationFS embed.FS

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open error: %s: %w", dsn, err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, fmt.Errorf("foreign keys pragma: %w", err)
	}

	if err := migrate(db, migrationFS, "*/*.sql"); err != nil {
		return nil, fmt.Errorf("migrate connect error: %w", err)
	}

	return db, nil
}

// migrate sets up migration tracking and executes pending migration files.
//
// Migration files are embedded in the sqlite/migration folder and are executed
// in lexigraphical order.
//
// Once a migration is run, its name is stored in the 'migrations' table so it
// is not re-executed. Migrations run in a transaction to prevent partial
// migrations.
func migrate(db *sql.DB, migrationFS fs.ReadFileFS, migrationsGlobPattern string) error {

	err := migrateMigrationTable(db)
	if err != nil {
		return fmt.Errorf("migrate migration table err: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, migrationsGlobPattern)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("no files found at: %v", migrationsGlobPattern)
	}
	sort.Strings(names)
	// fmt.Printf("migrations: %v\n", names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		if err := migrateFile(db, migrationFS, name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}
	return nil
}

// migrateMigrationTable creates the migrations table if it does not exist, and updates it if we have an old version of the table. This is a special table that tracks updates to the rest of the db, so we have to pudate it separately
func migrateMigrationTable(db *sql.DB) error {
	err := withTx(
		db,
		func(tx *sql.Tx) error {

			// Check if the migrations table exists
			migrationsTableCount := false
			if err := tx.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='migrations';`).Scan(&migrationsTableCount); err != nil {
				return fmt.Errorf("check if migrations table exists err: %w", err)
			}

			migrationV2TableCount := false
			if err := tx.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='migration_v2';`).Scan(&migrationV2TableCount); err != nil {
				return fmt.Errorf("check if migration_v2 table exists err: %w", err)
			}

			migrationV2Create := `
CREATE TABLE migration_v2 (
	migration_v2_id INTEGER PRIMARY KEY,
	file_name TEXT NOT NULL,
	migrate_time TEXT NOT NULL,
	UNIQUE(file_name)
) STRICT;
`

			switch {
			case migrationsTableCount && migrationV2TableCount:
				// this should never happen
				return errors.New("both migrations and migration_v2 tables exist")

			case migrationsTableCount && !migrationV2TableCount:
				if _, err := tx.Exec(migrationV2Create); err != nil {
					return fmt.Errorf("cannot create migration_v2 table: %w", err)
				}
				// NOTE: this relies on the embedded file system using '/' as the path separator. I think that's ok?
				insertIntoMigrationV2 := `
INSERT INTO migration_v2 (file_name, migrate_time)
SELECT replace(name, 'embedded_migrations/', '') AS file_name, ?
FROM migrations
ORDER BY name;
`
				now := time.Now().Round(0).UTC().Format(time.RFC3339)
				if _, err := tx.Exec(insertIntoMigrationV2, now); err != nil {
					return fmt.Errorf("cannot copy data from migrations to migration_v2: %w", err)
				}

				if _, err := tx.Exec(`DROP TABLE migrations;`); err != nil {
					return fmt.Errorf("cannot drop migrations table: %w", err)
				}

				return nil

			case !migrationsTableCount && migrationV2TableCount:
				return nil

			case !migrationsTableCount && !migrationV2TableCount:
				if _, err := tx.Exec(migrationV2Create); err != nil {
					return fmt.Errorf("cannot create migration_v2 table: %w", err)
				}

				return nil
			}

			return errors.New("unreachable")
		})

	if err != nil {
		return fmt.Errorf("migrate migration table err: %w", err)
	}

	return nil
}

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migrations" table to prevent re-running.
func migrateFile(db *sql.DB, migrationFS fs.ReadFileFS, name string) error {
	err := withTx(
		db,
		func(tx *sql.Tx) error {

			fileName := filepath.Base(name)
			// Ensure migration has not already been run.
			var n int
			if err := tx.QueryRow(`SELECT COUNT(*) FROM migration_v2 WHERE file_name = ?`, fileName).Scan(&n); err != nil {
				return err
			} else if n != 0 {
				return nil // already run migration, skip
			}

			// Read and execute migration file.
			if buf, err := fs.ReadFile(migrationFS, name); err != nil {
				return err
			} else if _, err := tx.Exec(string(buf)); err != nil {
				return err
			}

			now := time.Now().Round(0).UTC().Format(time.RFC3339)
			// Insert record into migrations to prevent re-running migration.
			if _, err := tx.Exec(`INSERT INTO migration_v2 (file_name, migrate_time) VALUES (?, ?)`, fileName, now); err != nil {
				return err
			}
			return nil
		},
	)

	if err != nil {
		return fmt.Errorf("migrate file err: %w", err)
	}
	return nil
}

// withTx makes transactions easy!!
func withTx(db *sql.DB, txFunc func(tx *sql.Tx) error) error {

	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("can't begin tx: %w", err)
		return err
	}
	// will not succeed if tx.Commit is called
	// explicitly ignore the error
	defer func() { _ = tx.Rollback() }()

	// do da magic
	err = txFunc(tx)

	if err != nil {
		err = fmt.Errorf("txFunc err: %w", err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit err: %w", err)
		return err
	}
	return nil
}
