package database

import (
	"context"
	"database/sql"

	// SQLite driver.
	_ "github.com/mattn/go-sqlite3"
)

// InitializeDB initializes the database and creates tables.
func InitializeDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if createErr := createUserTable(ctx, db); createErr != nil {
		_ = db.Close() // Ignore close error, return the original error
		return nil, createErr
	}

	return db, nil
}

func createUserTable(ctx context.Context, db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`

	_, err := db.ExecContext(ctx, query)
	return err
}
