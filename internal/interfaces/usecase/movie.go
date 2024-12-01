package usecase

import (
	"context"
	"lion-parcel-test/pkg/dto"
)

type MovieUsecase interface {
	CreateMovie(ctx context.Context, req *CreateMovieRequest) *dto.Response
	UpdateMovie(ctx context.Context, req *UpdateMovieRequest) *dto.Response
	MostViewed(ctx context.Context, req *MostViewedRequest) *dto.Response
	MostViewedGenre(ctx context.Context, req *MostViewedGenreRequest) *dto.Response
	GetMovies(ctx context.Context, req *GetMoviesRequest) *dto.Response
	SearchMovies(ctx context.Context, req *SearchMoviesRequest) *dto.Response
	VoteMovie(ctx context.Context, req *VoteMovieRequest) *dto.Response
	UnvoteMovie(ctx context.Context, req *UnvoteMovieRequest) *dto.Response
	VotedMovies(ctx context.Context, req *VotedMoviesRequest) *dto.Response
	MostVoted(ctx context.Context, req *MostVotedRequest) *dto.Response
	MostVotedGenre(ctx context.Context, req *MostVotedGenreRequest) *dto.Response
}

type VotedMoviesRequest struct {
	UserId int `json:"user_id" validate:"required"`
}
type VoteMovieRequest struct {
	UserId  int `json:"user_id" validate:"required"`
	MovieId int `json:"movie_id" validate:"required"`
}
type UnvoteMovieRequest struct {
	UserId  int `json:"user_id" validate:"required"`
	MovieId int `json:"movie_id" validate:"required"`
}

type SearchMoviesRequest struct {
	Keyword string
}
type SearchMoviesResponse struct {
	Movies interface{} `json:"movies"`
}
type VotedMoviesResponse struct {
	Movies interface{} `json:"movies"`
}

type GetMoviesRequest struct {
	Page     int
	PageSize int
}
type GetMoviesResponse struct {
	Movies         interface{} `json:"movies"`
	PaginationData interface{} `json:"pagination_data"`
}

type MostVotedGenreRequest struct {
}
type MostVotedRequest struct {
}
type MostViewedGenreRequest struct {
}
type MostViewedGenreResponse struct {
	ViewsCount int    `json:"views_count"`
	Genre      string `json:"genre"`
}
type MostVotedGenreResponse struct {
	VotesCount int    `json:"votes_count"`
	Genre      string `json:"genre"`
}
type MostViewedRequest struct {
}
type CreateMovieRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Duration    int    `json:"duration" validate:"required"`
	Artist      string `json:"artists" validate:"required"`
	Genre       string `json:"genres" validate:"required"`
	FileName    string `json:"file_name" validate:"required"`
}
type UpdateMovieRequest struct {
	Id          string `json:"id" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Duration    int    `json:"duration" validate:"required"`
	Artist      string `json:"artists" validate:"required"`
	Genre       string `json:"genres" validate:"required"`
	FileName    string `json:"file_name" validate:"required"`
}
