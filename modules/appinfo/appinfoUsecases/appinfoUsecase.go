package appinfoUsecases

import (
	"github.com/codepnw/ecommerce/modules/appinfo"
	"github.com/codepnw/ecommerce/modules/appinfo/appinfoRepositories"
)

type IAppinfoUsecase interface {
	FindCategory(req *appinfo.CategoryFilter) ([]*appinfo.Category, error)
	InsertCategory(req []*appinfo.Category) error
	DeleteCategory(categoryId int) error
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

func (u *appinfoUsecase) InsertCategory(req []*appinfo.Category) error {
	if err := u.repository.InsertCategory(req); err != nil {
		return err
	}
	return nil
}

func (u *appinfoUsecase) DeleteCategory(categoryId int) error {
	if err := u.repository.DeleteCategory(categoryId); err != nil {
		return err 
	}
	return nil
}