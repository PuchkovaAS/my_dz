package main

import (
	"4-order-api/configs"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf := configs.LoadConfig()

	dataBase, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	dataBase.AutoMigrate(&product.Product{}, &user.User{}, &order.Order{})
}
