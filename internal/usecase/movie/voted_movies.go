package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) VotedMovies(ctx context.Context, req *usecase.VotedMoviesRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "UnvoteMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	movie, err := uc.movieRepository.GetAllVotedMoviesByUserIdFromDb(ctx, req.UserId)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}
	resp.SetSuccess(http.StatusOK, "00", "Success get voted movie", usecase.VotedMoviesResponse{
		Movies: movie,
	})

	return resp
}
