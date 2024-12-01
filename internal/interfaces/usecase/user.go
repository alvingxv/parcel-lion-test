package usecase

import (
	"context"
	"lion-parcel-test/pkg/dto"
)

type UserUsecase interface {
	Login(ctx context.Context, req *LoginRequest) *dto.Response
	Register(ctx context.Context, req *RegisterRequest) *dto.Response
	PopulateSession(ctx context.Context, token string) *UserSession
}

type RegisterRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RegisterResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginRequest struct {
	Email string `json:"email"`
}
type LoginResponse struct {
	Jwt string `json:"jwt"`
}

type UserSession struct {
	Id      int  `json:"id"`
	IsAdmin bool `json:"is_admin"`
}
