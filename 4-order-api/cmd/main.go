package main

import (
	"4-order-api/configs"
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"4-order-api/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositorys
	productRepository := product.NewProductRepository(db)

	// Handlers
	product.NewProductHandler(
		router,
		product.ProductHandlerDeps{ProductRepository: productRepository},
	)

	// Middlewares
	stackMiddlewar := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stackMiddlewar(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
