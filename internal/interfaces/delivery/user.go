package delivery

import "github.com/gofiber/fiber/v2"

type UserHandler interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	IsAdmin(c *fiber.Ctx) error
	IsAuthenticated(c *fiber.Ctx) error
	PopulateSession(c *fiber.Ctx) error
}

type UserSession struct {
	Id      int  `json:"id"`
	IsAdmin bool `json:"is_admin"`
}
