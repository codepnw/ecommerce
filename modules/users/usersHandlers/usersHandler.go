package usersHandlers

import (
	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/modules/users"
	"github.com/codepnw/ecommerce/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type usersHandlersErrCode string

const (
	signUpCustomerErr usersHandlersErrCode = "users-001"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg     config.IConfig
	usecase usersUsecases.IUsersUsecase
}

func UsersHandler(cfg config.IConfig, usecase usersUsecases.IUsersUsecase) IUsersHandler {
	return &usersHandler{
		cfg:     cfg,
		usecase: usecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	req := new(users.UserRegisterReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}

	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email pattern is invalid",
		).Res()
	}

	result, err := h.usecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
		case "username has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}
