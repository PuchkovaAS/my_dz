package order

import (
	"4-order-api/internal/product"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	IsFormed bool           `json:"is_formed"`
	UserID   uint           `json:"user"      gorm:"index"`
	Items    []OrderProduct `json:"items"     gorm:"foreignKey:OrderID"`
}

type OrderProduct struct {
	OrderID   uint            `gorm:"primaryKey;autoIncrement:false"`
	ProductID uint            `gorm:"primaryKey;autoIncrement:false"`
	Product   product.Product `gorm:"foreignKey:ProductID"`
	Quantity  uint            `gorm:"default:1"`
}
