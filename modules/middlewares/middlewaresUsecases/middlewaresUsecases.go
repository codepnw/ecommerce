package middlewaresUsecases

import "github.com/codepnw/ecommerce/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecases interface {
	FindAccessToken(userId, accessToken string) bool
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