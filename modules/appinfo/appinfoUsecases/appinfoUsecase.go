package appinfoUsecases

import "github.com/codepnw/ecommerce/modules/appinfo/appinfoRepositories"

type IAppinfoUsecase interface {

}

type appinfoUsecase struct {
	repository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(repository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{repository: repository}
}