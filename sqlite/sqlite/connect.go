package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	_ "modernc.org/sqlite"
)

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
		return nil, fmt.Errorf("migrate: %w", err)
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
	// Ensure the 'migrations' table exists so we don't duplicate migrations.
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
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

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migrations" table to prevent re-running.
func migrateFile(db *sql.DB, migrationFS fs.ReadFileFS, name string) error {
	err := withTx(
		db,
		func(tx *sql.Tx) error {
			// Ensure migration has not already been run.
			var n int
			if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
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

			// Insert record into migrations to prevent re-running migration.
			if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
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
