package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) CreateMovie(ctx context.Context, req *usecase.CreateMovieRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "CreateMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	err := uc.movieRepository.InsertMovieToDB(ctx, req.Title, req.Description, req.Duration, req.Artist, req.Genre, req.FileName)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success Create Movie", nil)

	return resp
}
