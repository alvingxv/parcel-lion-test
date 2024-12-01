package useruc

import (
	"context"
	"errors"
	"lion-parcel-test/config"
	"lion-parcel-test/internal/interfaces/usecase"

	"github.com/golang-jwt/jwt"
	"go.elastic.co/apm/v2"
)

func (uc *userUsecase) PopulateSession(ctx context.Context, token string) *usecase.UserSession {
	apmSpan, ctx := apm.StartSpan(ctx, "IsAdmin", "usecase")
	defer apmSpan.End()

	tokenString, errp := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error parse")
		}
		return []byte(config.Cfg.Jwt.SecretKey), nil
	})
	if errp != nil {
		return nil
	}

	var mapClaims jwt.MapClaims

	if claims, ok := tokenString.Claims.(jwt.MapClaims); !ok || !tokenString.Valid {
		return nil
	} else {
		mapClaims = claims
	}

	var userEmail string
	if email, ok := mapClaims["email"].(string); !ok {
		return nil
	} else {
		userEmail = string(email)
	}
	user, err := uc.userRepository.GetUserFromDbByEmail(ctx, userEmail)
	if err != nil {
		return nil
	}

	return &usecase.UserSession{
		Id:      user.ID,
		IsAdmin: user.IsAdmin,
	}

}
