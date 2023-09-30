package middlewaresUsecases

import "github.com/codepnw/ecommerce/modules/middlewares/middlewaresRepositories"


type IMiddlewaresUsecases interface {
}

type middlewaresUsecases struct {
	repository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecases(repository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecases {
	return &middlewaresUsecases{
		repository: repository,
	}
}
