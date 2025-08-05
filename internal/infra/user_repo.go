package infra

import (
	"context"
	"database/sql"
	"errors"

	// SQLite driver.
	_ "github.com/mattn/go-sqlite3"

	"ministry-scheduler/internal/domain"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *SQLUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE email = ?`
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *SQLUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

func (r *SQLUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `UPDATE users SET name = ?, email = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, user.Name, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *SQLUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *SQLUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if scanErr := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); scanErr != nil {
			return nil, scanErr
		}
		users = append(users, &user)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return users, nil
}

func InitializeDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	if createErr := createTable(ctx, db); createErr != nil {
		_ = db.Close() // Ignore close error, return the original error
		return nil, createErr
	}

	return db, nil
}

func createTable(ctx context.Context, db *sql.DB) error {
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
