package useruc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *userUsecase) Register(ctx context.Context, req *usecase.RegisterRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "Register", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	err := uc.userRepository.InsertUserToDB(ctx, req.Email, req.Name)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success Register", nil)

	return resp
}
