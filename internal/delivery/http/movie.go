package http

import (
	"lion-parcel-test/constant"
	"lion-parcel-test/internal/interfaces/delivery"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.elastic.co/apm/v2"
)

type movieHandler struct {
	movieUsecase usecase.MovieUsecase
	validate     *validator.Validate
}

func NewMovieHandler(movieUsecase usecase.MovieUsecase, validate *validator.Validate) delivery.MovieHandler {
	return &movieHandler{
		movieUsecase: movieUsecase,
		validate:     validate,
	}
}

func (h *movieHandler) CreateMovie(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "CreateMovie", "Handler")
	defer apmSpan.End()

	jsonformfile := c.FormValue("json")

	var reqStruct usecase.CreateMovieRequest

	err := jsoniter.UnmarshalFromString(jsonformfile, &reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusUnprocessableEntity)
		c.JSON(dto.NewError(http.StatusUnprocessableEntity, "FM", "error unmarshall", err))
		return nil
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("File upload failed: " + err.Error())
	}
	reqStruct.FileName = strings.ReplaceAll(file.Filename, " ", "-")

	err = h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "FM", "error unmarshall", err))
		return nil
	}

	uploadDir := "./movies"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusInternalServerError)
		c.JSON(dto.NewError(http.StatusInternalServerError, "FM", "Failed to create upload directory", err))
		return nil
	}

	savePath := filepath.Join(uploadDir, reqStruct.FileName)
	if err := c.SaveFile(file, savePath); err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusInternalServerError)
		c.JSON(dto.NewError(http.StatusInternalServerError, "FM", "Failed to save movie file", err))
	}

	resp := h.movieUsecase.CreateMovie(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) UpdateMovie(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "CreateMovie", "Handler")
	defer apmSpan.End()

	jsonformfile := c.FormValue("json")

	var reqStruct usecase.UpdateMovieRequest

	err := jsoniter.UnmarshalFromString(jsonformfile, &reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusUnprocessableEntity)
		c.JSON(dto.NewError(http.StatusUnprocessableEntity, "FM", "error unmarshall", err))
		return nil
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("File upload failed: " + err.Error())
	}
	reqStruct.FileName = strings.ReplaceAll(file.Filename, " ", "-")

	id := c.Params("id")

	reqStruct.Id = id

	err = h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "FM", "error unmarshall", err))
		return nil
	}

	uploadDir := "./movies"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusInternalServerError)
		c.JSON(dto.NewError(http.StatusInternalServerError, "FM", "Failed to create upload directory", err))
		return nil
	}

	savePath := filepath.Join(uploadDir, reqStruct.FileName)
	if err := c.SaveFile(file, savePath); err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusInternalServerError)
		c.JSON(dto.NewError(http.StatusInternalServerError, "FM", "Failed to save movie file", err))
	}

	resp := h.movieUsecase.UpdateMovie(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) MostViewed(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "MostViewed", "Handler")
	defer apmSpan.End()

	resp := h.movieUsecase.MostViewed(ctx, nil)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) MostViewedGenre(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "MostViewed", "Handler")
	defer apmSpan.End()

	resp := h.movieUsecase.MostViewedGenre(ctx, nil)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) GetMovies(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "GetMovies", "Handler")
	defer apmSpan.End()

	page := c.Query("page", "1")
	pageSize := c.Query("pageSize", "10")

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)

	var reqStruct usecase.GetMoviesRequest

	reqStruct.Page = pageInt
	reqStruct.PageSize = pageSizeInt

	resp := h.movieUsecase.GetMovies(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) SearchMovies(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "SearchMovies", "Handler")
	defer apmSpan.End()

	keyword := c.Query("keyword", "")

	var reqStruct usecase.SearchMoviesRequest

	reqStruct.Keyword = keyword

	resp := h.movieUsecase.SearchMovies(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) VoteMovie(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "VoteMovie", "Handler")
	defer apmSpan.End()

	reqBody := c.Body()

	var reqStruct usecase.VoteMovieRequest

	err := jsoniter.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusUnprocessableEntity)
		c.JSON(dto.NewError(http.StatusBadRequest, "FM", "error unmarshall", err))
		return nil
	}

	session := c.Locals(constant.UserSessionKey).(*usecase.UserSession)

	reqStruct.UserId = session.Id

	err = h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "VE", "Validation Error", err))
		return nil
	}

	resp := h.movieUsecase.VoteMovie(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) UnvoteMovie(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "UnvoteMovie", "Handler")
	defer apmSpan.End()

	reqBody := c.Body()

	var reqStruct usecase.UnvoteMovieRequest

	err := jsoniter.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusUnprocessableEntity)
		c.JSON(dto.NewError(http.StatusBadRequest, "FM", "error unmarshall", err))
		return nil
	}

	session := c.Locals(constant.UserSessionKey).(*usecase.UserSession)

	reqStruct.UserId = session.Id

	err = h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "VE", "Validation Error", err))
		return nil
	}

	resp := h.movieUsecase.UnvoteMovie(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) VotedMovies(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "VotedMovies", "Handler")
	defer apmSpan.End()

	var reqStruct usecase.VotedMoviesRequest

	session := c.Locals(constant.UserSessionKey).(*usecase.UserSession)

	reqStruct.UserId = session.Id

	err := h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "VE", "Validation Error", err))
		return nil
	}

	resp := h.movieUsecase.VotedMovies(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) MostVoted(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "MostVoted", "Handler")
	defer apmSpan.End()

	resp := h.movieUsecase.MostVoted(ctx, nil)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *movieHandler) MostVotedGenre(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "MostVotedGenre", "Handler")
	defer apmSpan.End()

	resp := h.movieUsecase.MostVotedGenre(ctx, nil)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}
