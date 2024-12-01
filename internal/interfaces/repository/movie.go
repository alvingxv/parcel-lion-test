package repository

import (
	"context"
	"lion-parcel-test/pkg/errs"
)

type MovieRepository interface {
	InsertMovieToDB(ctx context.Context, Title string, Description string, Duration int, Artist string, Genre string, FileName string) errs.MessageErr
	UpdateMovieToDB(ctx context.Context, Id string, Title string, Description string, Duration int, Artist string, Genre string, FileName string) errs.MessageErr
	GetMostViewedMovieFromDB(ctx context.Context) (Movie, errs.MessageErr)
	GetMostViewedGenreFromDB(ctx context.Context) (Movie, errs.MessageErr)
	GetMoviesFromDB(ctx context.Context, page int, pageSize int) ([]Movie, MoviePaginationMetadata, errs.MessageErr)
	SearchMoviesFromDB(ctx context.Context, keyword string) ([]Movie, errs.MessageErr)
	InsertVoteToDB(ctx context.Context, userId int, movieId int) errs.MessageErr
	DeleteVoteFromDB(ctx context.Context, userId int, movieId int) errs.MessageErr
	GetAllVotedMoviesByUserIdFromDb(ctx context.Context, userId int) ([]Movie, errs.MessageErr)
	GetMostVotedMovieFromDB(ctx context.Context) (Movie, errs.MessageErr)
	GetMostVotedGenreFromDB(ctx context.Context) (Movie, errs.MessageErr)
	// GetUserFromDbByEmail(ctx context.Context, email string) (User, errs.MessageErr)
}

type Movie struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Artist      string `json:"artists"`
	Genre       string `json:"genres"`
	WatchUrl    string `json:"watch_url"`
	Views       int    `json:"views"`
	Vote        int    `json:"vote,omitempty"`
}
type MoviePaginationMetadata struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
}
