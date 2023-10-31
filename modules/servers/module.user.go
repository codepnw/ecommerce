package servers

import (
	"github.com/codepnw/ecommerce/modules/users/usersHandlers"
	"github.com/codepnw/ecommerce/modules/users/usersRepositories"
	"github.com/codepnw/ecommerce/modules/users/usersUsecases"
)

type IUsersModule interface {
	Init()
	Repository() usersRepositories.IUsersRepository
	Usecase() usersUsecases.IUsersUsecase
	Handler() usersHandlers.IUsersHandler
}

type usersModule struct {
	*moduleFactory
	repository usersRepositories.IUsersRepository
	usecase    usersUsecases.IUsersUsecase
	handler    usersHandlers.IUsersHandler
}

func (m *moduleFactory) UsersModule() IUsersModule {
	repository := usersRepositories.UsersRepository(m.s.db)
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repository)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	return &usersModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (u *usersModule) Init() {
	router := u.r.Group("/users")
	router.Get("/:user_id", u.m.JwtAuth(), u.m.ParamsCheck(), u.handler.GetUserProfile)
	router.Get("/admin/secret", u.m.JwtAuth(), u.m.Authorize(2), u.handler.GenerateAdminToken)

	router.Post("/signup", u.m.ApiKeyAuth(), u.handler.SignUpCustomer)
	router.Post("/signin", u.m.ApiKeyAuth(), u.handler.SignIn)
	router.Post("/refresh", u.m.ApiKeyAuth(), u.handler.RefreshPassport)
	router.Post("/signout", u.m.ApiKeyAuth(), u.handler.SignOut)
	router.Post("/signup-admin", u.m.JwtAuth(), u.m.Authorize(2), u.handler.SignOut)
}

func (u *usersModule) Repository() usersRepositories.IUsersRepository { return u.repository }
func (u *usersModule) Usecase() usersUsecases.IUsersUsecase           { return u.usecase }
func (u *usersModule) Handler() usersHandlers.IUsersHandler           { return u.handler }
