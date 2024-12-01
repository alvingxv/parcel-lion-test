package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) SearchMovies(ctx context.Context, req *usecase.SearchMoviesRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "CreateMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	movies, err := uc.movieRepository.SearchMoviesFromDB(ctx, req.Keyword)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success get movie", usecase.SearchMoviesResponse{
		Movies: movies,
	})

	return resp
}
