package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) UpdateMovie(ctx context.Context, req *usecase.UpdateMovieRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "UpdateMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	err := uc.movieRepository.UpdateMovieToDB(ctx, req.Id, req.Title, req.Description, req.Duration, req.Artist, req.Genre, req.FileName)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success Update Movie", nil)

	return resp
}
