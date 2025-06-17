package user

import (
	"4-order-api/internal/order"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone     string
	SessionId string
	Code      uint
	Orders    []order.Order `gorm:"foreignKey:UserID"`
}
