package main

import (
	"4-order-api/configs"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
)

func main() {
	conf := configs.LoadConfig()

	dataBase := db.NewDb(conf)
	dataBase.AutoMigrate(&product.Product{}, &user.UserId{})
}
