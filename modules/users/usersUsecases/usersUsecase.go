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
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
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

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	claims, err := auth.ParseToken(u.cfg.Jwt(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	oauth, err := u.repository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	profile, err := u.repository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id: profile.Id,
		RoleId: profile.RoleId,
	}

	accessToken, err := auth.NewEcomAuth(
		auth.Access,
		u.cfg.Jwt(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RepeatToken(
		u.cfg.Jwt(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)
	
	passport := &users.UserPassport{
		User: profile,
		Token: &users.UserToken{
			Id: oauth.Id,
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	if err := u.repository.UpdateOauth(passport.Token); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.repository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}