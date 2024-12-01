package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"lion-parcel-test/internal/interfaces/adapter"
	"lion-parcel-test/internal/interfaces/repository"
	"lion-parcel-test/pkg/errs"

	"go.elastic.co/apm/v2"
)

type userRepository struct {
	database adapter.DatabaseClient
}

func NewUserRepository(database adapter.DatabaseClient) repository.UserRepository {
	return &userRepository{
		database: database,
	}
}

func (rp *userRepository) InsertUserToDB(ctx context.Context, email string, name string) errs.MessageErr {
	apmSpan, ctx := apm.StartSpan(ctx, "InsertUserToDB", "Repository")
	defer apmSpan.End()

	query := `INSERT INTO users (name, email, is_admin) VALUES (?, ?, false);`

	result := rp.database.Execute(ctx, query, name, email)
	if result.Error != nil {
		return errs.NewCustomErrs(
			"Failed Insert Database",
			"FD",
			result.Error.Error(),
		)
	}

	return nil
}
func (rp *userRepository) GetUserFromDbByEmail(ctx context.Context, email string) (repository.User, errs.MessageErr) {
	apmSpan, ctx := apm.StartSpan(ctx, "GetUserFromDbByEmail", "Repository")
	defer apmSpan.End()

	query := `SELECT * FROM users WHERE email = ?`

	var user repository.User

	row := rp.database.QueryRow(ctx, query, email)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.User{}, errs.NewCustomErrs(
				"Not Exist",
				"NA",
				err.Error(),
			)
		}

		return repository.User{}, errs.NewCustomErrs(
			"Failed Get Database",
			"FD",
			err.Error(),
		)
	}

	return user, nil
}
