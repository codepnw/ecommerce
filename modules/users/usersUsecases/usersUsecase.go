package usersUsecases

import (
	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/users"
	"github.com/codepnw/ecommerce/modules/users/usersRepositories"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg config.IConfig
	repository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, repository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg: cfg,
		repository: repository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	result, err := u.repository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}
	return result, nil
}