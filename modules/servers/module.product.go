package servers

import (
	"github.com/codepnw/ecommerce/modules/products/productsHandlers"
	"github.com/codepnw/ecommerce/modules/products/productsRepositories"
	"github.com/codepnw/ecommerce/modules/products/productsUsecases"
)

type IProductsModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type productsModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductsModule {
	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(m.s.cfg, usecase, m.FilesModule().Usecase())

	return &productsModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (p *productsModule) Init() {
	router := p.r.Group("/products")
	router.Post("/", p.m.JwtAuth(), p.m.Authorize(2), p.handler.InsertProduct)
	router.Patch("/:product_id", p.m.JwtAuth(), p.m.Authorize(2), p.handler.UpdateProduct)

	router.Get("/", p.m.ApiKeyAuth(), p.handler.FindProduct)
	router.Get("/:product_id", p.m.ApiKeyAuth(), p.handler.FindOneProduct)

	router.Delete("/:product_id", p.m.JwtAuth(), p.m.Authorize(2), p.handler.DeleteProduct)
}

func (p *productsModule) Repository() productsRepositories.IProductsRepository { return p.repository }
func (p *productsModule) Usecase() productsUsecases.IProductsUsecase           { return p.usecase }
func (p *productsModule) Handler() productsHandlers.IProductsHandler           { return p.handler }
