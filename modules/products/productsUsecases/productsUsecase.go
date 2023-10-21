package productsUsecases

import (
	"math"

	"github.com/codepnw/ecommerce/modules/entities"
	"github.com/codepnw/ecommerce/modules/products"
	"github.com/codepnw/ecommerce/modules/products/productsRepositories"
)

type IProductsUsecase interface {
	FindOneProduct(productId string) (*products.Product, error)
	FindProduct(req *products.ProductFilter) *entities.PaginateRes
	InsertProduct(req *products.Product) (*products.Product, error)
	DeleteProduct(productId string) error
	UpdateProduct(req *products.Product) (*products.Product, error) 
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

func (u *productsUsecase) FindProduct(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.repository.FindProduct(req)

	return &entities.PaginateRes{
		Data: products,
		Page: req.Page,
		Limit: req.Limit,
		TotalItem: count,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}
}

func (u *productsUsecase) InsertProduct(req *products.Product) (*products.Product, error) {
	product, err := u.repository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecase) DeleteProduct(productId string) error {
	if err := u.repository.DeleteProduct(productId); err != nil {
		return err
	}
	return nil
}

func (u *productsUsecase) UpdateProduct(req *products.Product) (*products.Product, error) {
	product, err := u.repository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}