package order

import (
	"4-order-api/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	IsFormed bool               `json:"is_formed"`
	UserID   uint               `json:"user"      gorm:"index"`
	Products []*product.Product `json:"products"  gorm:"many2many:order_products;"`
	Items    []OrderProduct     `json:"items"     gorm:"foreignKey:OrderID"`
}

type OrderProduct struct {
	OrderID   uint            `gorm:"primaryKey"`
	ProductID uint            `gorm:"primaryKey"`
	Quantity  uint            `gorm:"default:1"`
	Product   product.Product `gorm:"foreignKey:ProductID"`
}
