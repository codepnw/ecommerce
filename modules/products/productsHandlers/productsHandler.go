package productsHandlers

import (
	"strings"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/modules/files/filesUsecases"
	"github.com/codepnw/ecommerce/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHnadlerErrCode string

const (
	findOneProductErr productsHnadlerErrCode = "products-001"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
}

type productsHnalder struct {
	cfg config.IConfig	
	usecase productsUsecases.IProductsUsecase
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.IConfig, usecase productsUsecases.IProductsUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHnalder{
		cfg: cfg,
		usecase: usecase,
		filesUsecase: filesUsecase,
	}
}

func (h *productsHnalder) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.usecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}
