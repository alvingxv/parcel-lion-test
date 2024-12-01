package app

import (
	"lion-parcel-test/internal/interfaces/usecase"
	movieuc "lion-parcel-test/internal/usecase/movie"
	useruc "lion-parcel-test/internal/usecase/user"
)

type Usecases struct {
	UserUsecase  usecase.UserUsecase
	MovieUsecase usecase.MovieUsecase
}

func NewUsecases(repos *Repositories) *Usecases {

	return &Usecases{
		UserUsecase:  useruc.NewUserUsecase(repos.userRepository),
		MovieUsecase: movieuc.NewMovieUsecase(repos.movieRepository),
	}
}
