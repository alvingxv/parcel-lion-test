package movieuc

import (
	"lion-parcel-test/internal/interfaces/repository"
	"lion-parcel-test/internal/interfaces/usecase"
)

type movieUsecase struct {
	movieRepository repository.MovieRepository
}

func NewMovieUsecase(movieRepository repository.MovieRepository) usecase.MovieUsecase {
	return &movieUsecase{
		movieRepository: movieRepository,
	}
}
