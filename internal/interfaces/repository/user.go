package repository

import (
	"context"
	"lion-parcel-test/pkg/errs"
	"time"
)

type UserRepository interface {
	InsertUserToDB(ctx context.Context, email string, name string) errs.MessageErr
	GetUserFromDbByEmail(ctx context.Context, email string) (User, errs.MessageErr)
}

type User struct {
	ID        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	IsAdmin   bool      `db:"is_admin" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
