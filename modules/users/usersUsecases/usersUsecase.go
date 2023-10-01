package usersUsecases

import (
	"fmt"

	"github.com/codepnw/ecommerce/config"
	"github.com/codepnw/ecommerce/modules/users"
	"github.com/codepnw/ecommerce/modules/users/usersRepositories"
	"github.com/codepnw/ecommerce/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
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

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.repository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	accessToken, err := auth.NewEcomAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})

	refreshToken, err := auth.NewEcomAuth(auth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})

	passport := &users.UserPassport{
		User: &users.User{
			Id: user.Id,
			Email: user.Email,
			Username: user.Username,
			RoleId: user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}

	if err := u.repository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}



