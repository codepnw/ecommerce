package appinfoHandlers

import (
	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/appinfo"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoUsecases"
	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlerErrCode string

const (
	generateApiKeyErr appinfoHandlerErrCode = "appinfo-001"
	findCategoryErr appinfoHandlerErrCode = "appinfo-002"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error	
}

type appinfoHandler struct {
	cfg config.IConfig
	usecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, usecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg: cfg,
		usecase: usecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewEcomAuth(
		auth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)

	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		} {
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.usecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}