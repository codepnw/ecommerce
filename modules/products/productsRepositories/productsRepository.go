package productsRepositories

type IProductsRepository interface {

}

type productRepository struct {

}

func ProductsRepository() IProductsRepository {
	return &productRepository{}
}