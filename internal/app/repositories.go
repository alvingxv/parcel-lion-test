package app

import (
	"lion-parcel-test/internal/interfaces/repository"
	movierepo "lion-parcel-test/internal/repository/movie"
	userrepo "lion-parcel-test/internal/repository/user"
)

type Repositories struct {
	userRepository  repository.UserRepository
	movieRepository repository.MovieRepository
}

func NewRepos(dependencies *Dependencies) *Repositories {
	return &Repositories{
		userRepository:  userrepo.NewUserRepository(dependencies.sqlitedb),
		movieRepository: movierepo.NewMovieRepository(dependencies.sqlitedb),
	}
}
