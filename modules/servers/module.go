package servers

import (
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresHandlers"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresRepositories"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresUsecases"
	monitorHandlers "github.com/codepnw/ecommerce/modules/monitor/handlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule() IUsersModule
	AppinfoModule() IAppinfoModule
	FilesModule() IFilesModule
	ProductsModule() IProductsModule
	OrdersModule() IOrdersModule
}

type moduleFactory struct {
	r fiber.Router
	s *server
	m middlewaresHandlers.IMiddlewaresHandlers
}

func InitModule(r fiber.Router, s *server, m middlewaresHandlers.IMiddlewaresHandlers) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
		m: m,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandlers {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecases(repository)
	return middlewaresHandlers.MiddlewaresHandlers(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.s.cfg)
	m.r.Get("/", handler.HealthCheck)
}

