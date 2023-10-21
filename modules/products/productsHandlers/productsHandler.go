package productsHandlers

import (
	"strings"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/appinfo"
	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/modules/files/filesUsecases"
	"github.com/codepnw/ecommerce/modules/products"
	"github.com/codepnw/ecommerce/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHnadlerErrCode string

const (
	findOneProductErr productsHnadlerErrCode = "products-001"
	findProductErr    productsHnadlerErrCode = "products-002"
	insertProductErr    productsHnadlerErrCode = "products-003"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	InsertProduct(c *fiber.Ctx) error
}

type productsHnalder struct {
	cfg          config.IConfig
	usecase      productsUsecases.IProductsUsecase
	filesUsecase filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.IConfig, usecase productsUsecases.IProductsUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHnalder{
		cfg:          cfg,
		usecase:      usecase,
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

func (h *productsHnalder) FindProduct(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.usecase.FindProduct(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHnalder) InsertProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images: make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.usecase.InsertProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}