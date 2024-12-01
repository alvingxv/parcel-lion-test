package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"lion-parcel-test/constant"
	"lion-parcel-test/internal/interfaces/adapter"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"go.elastic.co/apm/v2"
)

type sqliteClient struct {
	db *sql.DB
}

func NewSqliteClient() (adapter.DatabaseClient, error) {
	db, err := sql.Open("sqlite3", "./movies.db")
	if err != nil {
		return nil, err
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createUsersTable)
	if err != nil {
		return nil, err
	}

	createMoviesTable := `CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL,
    artists TEXT,
    genres TEXT,
    watch_url TEXT,
    views_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createMoviesTable)
	if err != nil {
		return nil, err
	}

	createVotesTable := `CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    movie_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, movie_id),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(movie_id) REFERENCES movies(id) ON DELETE CASCADE
	);`
	_, err = db.Exec(createVotesTable)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &sqliteClient{
		db: db,
	}, nil
}

func (r *sqliteClient) Execute(ctx context.Context, query string, args ...interface{}) adapter.ExecuteResult {
	span, ctx := apm.StartSpan(ctx, "Execute", "database")
	defer span.End()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return adapter.ExecuteResult{
				Error: errors.New(constant.DuplicateConstraintError),
			}
		}

		apm.CaptureError(ctx, err)
		return adapter.ExecuteResult{
			Error: fmt.Errorf("failed to execute query: %w", err),
		}
	}

	lastInsertID, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	return adapter.ExecuteResult{
		Result:       result,
		LastInsertID: lastInsertID,
		RowsAffected: rowsAffected,
		Error:        nil,
	}
}

func (r *sqliteClient) QueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := apm.StartSpan(ctx, "QueryRows", "database")
	defer span.End()

	// Execute the query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		apm.CaptureError(ctx, err)
		return nil, err
	}

	return rows, nil
}

func (r *sqliteClient) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	span, ctx := apm.StartSpan(ctx, "QueryRow", "database")
	defer span.End()

	// Execute the query
	row := r.db.QueryRowContext(ctx, query, args...)

	return row
}

func (s *sqliteClient) Close() error {
	err := s.db.Close()

	if err != nil {
		return err
	}

	return nil
}
