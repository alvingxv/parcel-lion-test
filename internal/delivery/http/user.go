package http

import (
	"lion-parcel-test/constant"
	"lion-parcel-test/internal/interfaces/delivery"
	"lion-parcel-test/internal/interfaces/usecase"
	"lion-parcel-test/pkg/dto"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"go.elastic.co/apm/v2"
)

type userHandler struct {
	userUsecase usecase.UserUsecase
	validate    *validator.Validate
}

func NewUserHandler(userUsecase usecase.UserUsecase, validate *validator.Validate) delivery.UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
		validate:    validate,
	}
}

func (h *userHandler) Register(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "Register", "Handler")
	defer apmSpan.End()

	reqBody := c.Body()

	var reqStruct usecase.RegisterRequest

	err := jsoniter.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusUnprocessableEntity)
		c.JSON(dto.NewError(http.StatusUnprocessableEntity, "FM", "error unmarshall", err))
		return nil
	}

	err = h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "FM", "error unmarshall", err))
		return nil
	}

	resp := h.userUsecase.Register(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "Register", "Handler")
	defer apmSpan.End()

	reqBody := c.Body()

	var reqStruct usecase.LoginRequest

	err := jsoniter.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusUnprocessableEntity)
		c.JSON(dto.NewError(http.StatusBadRequest, "FM", "error unmarshall", err))
		return nil
	}

	err = h.validate.Struct(reqStruct)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		c.Status(http.StatusBadRequest)
		c.JSON(dto.NewError(http.StatusBadRequest, "VE", "Validation Error", err))
		return nil
	}

	resp := h.userUsecase.Login(ctx, &reqStruct)

	c.Status(resp.HttpCode)
	c.JSON(resp)
	return nil
}

func (h *userHandler) PopulateSession(c *fiber.Ctx) error {
	apmSpan, ctx := apm.StartSpan(c.Context(), "IsAdmin", "Handler")
	defer apmSpan.End()

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		c.Next()
		return nil
	}

	isBearer := strings.HasPrefix(authHeader, "Bearer")

	if !isBearer {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

	splitToken := strings.Split(authHeader, " ")

	if len(splitToken) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

	tokenString := splitToken[1]

	userSession := h.userUsecase.PopulateSession(ctx, tokenString)

	if userSession == nil {
		c.Next()
		return nil
	}

	c.Locals(constant.UserSessionKey, userSession)
	c.Next()

	return nil
}

func (h *userHandler) IsAdmin(c *fiber.Ctx) error {
	apmSpan, _ := apm.StartSpan(c.Context(), "IsAdmin", "Handler")
	defer apmSpan.End()

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

	session := c.Locals(constant.UserSessionKey).(*usecase.UserSession)

	if session == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthenticated",
		})
	}

	if !session.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	c.Next()

	return nil
}

func (h *userHandler) IsAuthenticated(c *fiber.Ctx) error {
	apmSpan, _ := apm.StartSpan(c.Context(), "IsAdmin", "Handler")
	defer apmSpan.End()

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header is missing",
		})
	}

	session, ok := c.Locals(constant.UserSessionKey).(*usecase.UserSession)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthenticated",
		})
	}

	if session.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized You're admin",
		})
	}

	c.Next()

	return nil
}
