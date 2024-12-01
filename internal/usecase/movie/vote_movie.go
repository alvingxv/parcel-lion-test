package movieuc

import (
	"context"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"

	"go.elastic.co/apm/v2"
)

func (uc *movieUsecase) VoteMovie(ctx context.Context, req *usecase.VoteMovieRequest) *dto.Response {
	apmSpan, ctx := apm.StartSpan(ctx, "VoteMovie", "usecase")
	defer apmSpan.End()

	resp := dto.New()

	err := uc.movieRepository.InsertVoteToDB(ctx, req.UserId, req.MovieId)
	if err != nil {
		resp.SetError(http.StatusNotFound, err.Status(), err.Message(), err)
		return resp
	}

	resp.SetSuccess(http.StatusOK, "00", "Success vote", nil)

	return resp
}
