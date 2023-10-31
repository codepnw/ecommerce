package servers

import (
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoHandlers"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoRepositories"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoUsecases"
)

type IAppinfoModule interface {
	Init()
	Repository() appinfoRepositories.IAppinfoRepository
	Usecase() appinfoUsecases.IAppinfoUsecase
	Handler() appinfoHandlers.IAppinfoHandler
}

type appinfoModule struct {
	*moduleFactory
	repository appinfoRepositories.IAppinfoRepository
	usecase    appinfoUsecases.IAppinfoUsecase
	handler    appinfoHandlers.IAppinfoHandler
}

func (m *moduleFactory) AppinfoModule() IAppinfoModule {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.s.cfg, usecase)

	return &appinfoModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (a *appinfoModule) Init() {
	router := a.r.Group("/appinfo")
	router.Post("/categories", a.m.JwtAuth(), a.m.Authorize(2), a.handler.InsertCategory)

	router.Get("/apikey", a.m.JwtAuth(), a.m.Authorize(2), a.handler.GenerateApiKey)
	router.Get("/categories", a.m.ApiKeyAuth(), a.handler.FindCategory)

	router.Delete("/:category_id/categories", a.m.JwtAuth(), a.m.Authorize(2), a.handler.DeleteCategory)
}

func (a *appinfoModule) Repository() appinfoRepositories.IAppinfoRepository { return a.repository }
func (a *appinfoModule) Usecase() appinfoUsecases.IAppinfoUsecase           { return a.usecase }
func (a *appinfoModule) Handler() appinfoHandlers.IAppinfoHandler           { return a.handler }
