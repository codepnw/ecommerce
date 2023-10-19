package products

import (
	"github.com/codepnw/ecommerce/modules/appinfo"
	"github.com/codepnw/ecommerce/modules/entities"
)

type Product struct {
	Id          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Category    *appinfo.Category  `json:"category"`
	CreatedAt   string             `json:"created_at"`
	UpdatedAt   string             `json:"updated_at"`
	Price       float64            `json:"price"`
	Images      []*entities.Image `json:"images"`
}
