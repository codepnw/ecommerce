package middlewaresHandlers

import (
	"strings"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresUsecases"
	"github.com/codepnw/ecommerce/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewaresErrCode string

const (
	routerCheckErr middlewaresErrCode = "middleware-001"
	jwtAuthErr     middlewaresErrCode = "middleware-002"
	paramsCheckErr middlewaresErrCode = "middleware-003"
)

type IMiddlewaresHandlers interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
}

type middlewaresHandlers struct {
	cfg     config.IConfig
	usecase middlewaresUsecases.IMiddlewaresUsecases
}

func MiddlewaresHandlers(cfg config.IConfig, usecase middlewaresUsecases.IMiddlewaresUsecases) IMiddlewaresHandlers {
	return &middlewaresHandlers{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *middlewaresHandlers) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandlers) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandlers) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Asia/Bangkok",
	})
}

func (h *middlewaresHandlers) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.usecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandlers) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}
