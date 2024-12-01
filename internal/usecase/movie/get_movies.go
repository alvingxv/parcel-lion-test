package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) GetMovies(ctx context.Context, req *usecase.GetMoviesRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "CreateMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	movie, paginationMetadata, err := uc.movieRepository.GetMoviesFromDB(ctx, req.Page, req.PageSize)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success get movie", usecase.GetMoviesResponse{
		Movies:         movie,
		PaginationData: paginationMetadata,
	})

	return resp
}
