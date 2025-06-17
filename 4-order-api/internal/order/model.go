package order

import (
	"4-order-api/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID   uint              `json:"user"     gorm:"index"`
	Products []product.Product `json:"products" gorm:"many2many:order_products;"`
}
