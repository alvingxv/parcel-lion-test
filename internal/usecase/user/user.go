package useruc

import (
	"lion-parcel-test/internal/interfaces/repository"
	"lion-parcel-test/internal/interfaces/usecase"
)

type userUsecase struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) usecase.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}
