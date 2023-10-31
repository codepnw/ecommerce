package servers

import (
	"github.com/codepnw/ecommerce/modules/orders/ordersHandlers"
	"github.com/codepnw/ecommerce/modules/orders/ordersRepositories"
	"github.com/codepnw/ecommerce/modules/orders/ordersUsecases"
)

type IOrdersModule interface {
	Init()
	Repository() ordersRepositories.IOrdersRepository
	Usecase() ordersUsecases.IOrdersUsecase
	Handler() ordersHandlers.IOrdersHandler
}

type ordersModule struct {
	*moduleFactory
	repository ordersRepositories.IOrdersRepository
	usecase    ordersUsecases.IOrdersUsecase
	handler    ordersHandlers.IOrdersHandler
}

func (m *moduleFactory) OrdersModule() IOrdersModule {
	repository := ordersRepositories.OrdersRepository(m.s.db)
	usecase := ordersUsecases.OrdersUsecase(repository, m.ProductsModule().Repository())
	handler := ordersHandlers.OrdersHandler(m.s.cfg, usecase)

	return &ordersModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (o *ordersModule) Init() {
	router := o.r.Group("/orders")
	router.Post("/", o.m.JwtAuth(), o.handler.InsertOrder)

	router.Get("/", o.m.JwtAuth(), o.m.Authorize(2), o.handler.FindOrder)
	router.Get("/:user_id/:order_id", o.m.JwtAuth(), o.m.ParamsCheck(), o.handler.FindOneOrder)

	router.Patch("/:user_id/:order_id", o.m.JwtAuth(), o.m.ParamsCheck(), o.handler.UpdateOrder)
}

func (o *ordersModule) Repository() ordersRepositories.IOrdersRepository { return o.repository }
func (o *ordersModule) Usecase() ordersUsecases.IOrdersUsecase           { return o.usecase }
func (o *ordersModule) Handler() ordersHandlers.IOrdersHandler           { return o.handler }
