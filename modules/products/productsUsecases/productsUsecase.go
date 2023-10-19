package productsUsecases

import (
	"github.com/codepnw/ecommerce/modules/products"
	"github.com/codepnw/ecommerce/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
}

type productsUsecase struct {
	repository productsRepositories.IProductsRepository
}

func ProductsUsecase(repository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{repository: repository}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Product, error) {
	product, err := u.repository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}