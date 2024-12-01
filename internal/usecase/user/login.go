package useruc

import (
	"context"
	"lion-parcel-test/config"
	"lion-parcel-test/internal/interfaces/repository"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"go.elastic.co/apm/v2"
)

func (uc *userUsecase) Login(ctx context.Context, req *usecase.LoginRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "Register", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	user, err := uc.userRepository.GetUserFromDbByEmail(ctx, req.Email)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}

	token := generateToken(user)

	resp.SetSuccess(http.StatusOK, "00", "Success Login", usecase.LoginResponse{
		Jwt: token,
	})

	return resp
}

func generateToken(u repository.User) string {
	claims := jwt.MapClaims{
		"email": u.Email,
		"exp":   time.Now().Add(time.Hour * 10).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte(config.Cfg.Jwt.SecretKey))

	return tokenString
}
