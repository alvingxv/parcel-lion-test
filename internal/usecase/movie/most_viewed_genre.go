package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) MostViewedGenre(ctx context.Context, req *usecase.MostViewedGenreRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "CreateMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	movie, err := uc.movieRepository.GetMostViewedGenreFromDB(ctx)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success get most viewed genre movie", usecase.MostViewedGenreResponse{
		ViewsCount: movie.Views,
		Genre:      movie.Genre,
	})

	return resp
}
