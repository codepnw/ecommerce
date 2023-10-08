package middlewaresUsecases

import (
	"github.com/codepnw/ecommerce/modules/middlewares"
	"github.com/codepnw/ecommerce/modules/middlewares/middlewaresRepositories"
)

type IMiddlewaresUsecases interface {
	FindAccessToken(userId, accessToken string) bool
	FindRole() ([]*middlewares.Role, error)
}

type middlewaresUsecases struct {
	repository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecases(repository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		repository: repository,
	}
}

func (u *middlewaresUsecases) FindAccessToken(userId, accessToken string) bool {
	return u.repository.FindAccessToken(userId, accessToken)
}

func (u *middlewaresUsecases) FindRole() ([]*middlewares.Role, error) {
	roles, err := u.repository.FindRole()
	if err != nil {
		return nil, err
	}
	return roles, nil
}
