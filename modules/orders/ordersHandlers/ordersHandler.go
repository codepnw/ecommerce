package ordersHandlers

import (
	"strings"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlersErrCode string

const (
	findOneOrderErr ordersHandlersErrCode = "orders-001"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	usecase ordersUsecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, usecase ordersUsecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		usecase: usecase,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.usecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
