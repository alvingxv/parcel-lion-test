package adapter

import (
	"context"
	"database/sql"
)

type DatabaseClient interface {
	Close() error
	Execute(ctx context.Context, query string, args ...interface{}) ExecuteResult
	QueryRows(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type ExecuteResult struct {
	Result       sql.Result
	LastInsertID int64
	RowsAffected int64
	Error        error
}
