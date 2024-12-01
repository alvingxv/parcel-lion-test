package http

import (
	"context"
	"lion-parcel-test/config"
	"lion-parcel-test/constant"
	"lion-parcel-test/internal/app"
	"lion-parcel-test/pkg/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.elastic.co/apm/module/apmfiber/v2"
)

type HttpServer struct {
	*fiber.App
}

func NewHttpServer(app *app.App) (*HttpServer, error) {
	validate := validator.New()
	r := fiber.New()

	userHandler := NewUserHandler(app.Usecases.UserUsecase, validate)
	movieHandler := NewMovieHandler(app.Usecases.MovieUsecase, validate)

	r.Use(apmfiber.Middleware())
	r.Use(middleware.LoggingMiddleware)
	r.Use(userHandler.PopulateSession)

	// All users
	r.Post(constant.RouteApiV1+"/register", userHandler.Register)
	r.Post(constant.RouteApiV1+"/login", userHandler.Login)

	r.Get(constant.RouteApiV1+"/movies", movieHandler.GetMovies)
	r.Get(constant.RouteApiV1+"/movies/search", movieHandler.SearchMovies)

	// admin
	adminR := r.Group(constant.RouteApiV1+"/admin", userHandler.IsAdmin)
	adminR.Post("/movies", movieHandler.CreateMovie)
	adminR.Put("/movies/:id", movieHandler.UpdateMovie)
	adminR.Get("/movies/most_viewed", movieHandler.MostViewed)
	adminR.Get("/movies/most_viewed_genre", movieHandler.MostViewedGenre)
	adminR.Get("/movies/most_voted", movieHandler.MostVoted)
	adminR.Get("/movies/most_voted_genre", movieHandler.MostVotedGenre)

	// authenticated user
	authUser := r.Group(constant.RouteApiV1+"/movies", userHandler.IsAuthenticated)
	authUser.Post("/vote", movieHandler.VoteMovie)
	authUser.Post("/unvote", movieHandler.UnvoteMovie)
	authUser.Get("/votes", movieHandler.VotedMovies)

	// middeware to add view count of that movies
	// app.Use("/uploads", staticFileMiddleware)
	r.Static("/movies", "./movies")

	r.Get("/healthz", func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Status(200)
		return c.SendString("ok")
	})

	return &HttpServer{
		r,
	}, nil
}

func (s *HttpServer) Run() error {
	return s.Listen(":" + config.Cfg.App.Port)
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.Shutdown()
}
