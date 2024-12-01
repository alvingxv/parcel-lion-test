package delivery

import "github.com/gofiber/fiber/v2"

type MovieHandler interface {
	CreateMovie(c *fiber.Ctx) error
	UpdateMovie(c *fiber.Ctx) error
	MostViewed(c *fiber.Ctx) error
	MostViewedGenre(c *fiber.Ctx) error
	GetMovies(c *fiber.Ctx) error
	SearchMovies(c *fiber.Ctx) error
	VoteMovie(c *fiber.Ctx) error
	UnvoteMovie(c *fiber.Ctx) error
	VotedMovies(c *fiber.Ctx) error
	MostVoted(c *fiber.Ctx) error
	MostVotedGenre(c *fiber.Ctx) error
}
