package appinfoUsecases

import (
	"github.com/codepnw/ecommerce/modules/appinfo"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
}

type appinfoUsecase struct {
	repository appinfoRepositories.IAppinfoRepository
}

func AppinfoUsecase(repository appinfoRepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{repository: repository}
}

func (u *appinfoUsecase) FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error) {
	category, err := u.repository.FindCategory(req)
	if err != nil {
		return nil, err
	}
	return category, nil
}