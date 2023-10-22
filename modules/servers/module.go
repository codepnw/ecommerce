package servers

import (
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoHandlers"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoRepositories"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoUsecases"
	"github.com/codepnw/ecommerce/modules/files/filesHandlers"
	"github.com/codepnw/ecommerce/modules/files/filesUsecases"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresHandlers"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresRepositories"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresUsecases"
	monitorHandlers "github.com/codepnw/ecommerce/modules/monitor/handlers"
	"github.com/codepnw/ecommerce/modules/orders/ordersHandlers"
	"github.com/codepnw/ecommerce/modules/orders/ordersRepositories"
	"github.com/codepnw/ecommerce/modules/orders/ordersUsecases"
	"github.com/codepnw/ecommerce/modules/products/productsHandlers"
	"github.com/codepnw/ecommerce/modules/products/productsRepositories"
	"github.com/codepnw/ecommerce/modules/products/productsUsecases"
	"github.com/codepnw/ecommerce/modules/users/usersHandlers"
	"github.com/codepnw/ecommerce/modules/users/usersRepositories"
	"github.com/codepnw/ecommerce/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
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

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")

	router.Post("/signup", m.m.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.m.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.m.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.m.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.m.JwtAuth(), m.m.Authorize(2), handler.SignOut)

	router.Get("/:user_id", m.m.JwtAuth(), m.m.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.m.JwtAuth(), m.m.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	router := m.r.Group("/appinfo")

	router.Post("/categories", m.m.JwtAuth(), m.m.Authorize(2), handler.InsertCategory)

	router.Get("/apikey", m.m.JwtAuth(), m.m.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.m.ApiKeyAuth(), handler.FindCategory)

	router.Delete("/:category_id/categories", m.m.JwtAuth(), m.m.Authorize(2), handler.DeleteCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FilesHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.m.JwtAuth(), m.m.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.m.JwtAuth(), m.m.Authorize(2), handler.DeleteFile)
}

func (m *moduleFactory) ProductsModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)

	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)
	usecase := productsUsecases.ProductsUsecase(repository)
	handlers := productsHandlers.ProductsHandler(m.s.cfg, usecase, filesUsecase)

	router := m.r.Group("/products")

	router.Post("/", m.m.JwtAuth(), m.m.Authorize(2), handlers.InsertProduct)
	router.Patch("/:product_id", m.m.JwtAuth(), m.m.Authorize(2), handlers.UpdateProduct)

	router.Get("/", m.m.ApiKeyAuth(), handlers.FindProduct)
	router.Get("/:product_id", m.m.ApiKeyAuth(), handlers.FindOneProduct)

	router.Delete("/:product_id", m.m.JwtAuth(), m.m.Authorize(2), handlers.DeleteProduct)

}

func (m *moduleFactory) OrdersModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productsRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, filesUsecase)

	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	usecase := ordersUsecases.OrdersUsecase(ordersRepository, productsRepository)
	handler := ordersHandlers.OrdersHandler(m.s.cfg, usecase)

	router := m.r.Group("/orders")

	router.Get("/:order_id", m.m.JwtAuth(), m.m.Authorize(2), handler.FindOneOrder)
}
