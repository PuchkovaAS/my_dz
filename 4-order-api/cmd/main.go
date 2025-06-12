package main

import (
	"4-order-api/configs"
	"4-order-api/internal/auth"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Logger
	logger := middleware.NewLogger(conf)

	// JWT
	jwtService := jwt.NewJWT(conf.Auth.Secret)

	// Repositorys
	productRepository := product.NewProductRepository(db)
	userRepository := user.NewUserRepository(db)

	// Services
	authService := auth.NewUserService(userRepository)

	// Handlers
	product.NewProductHandler(
		router,
		product.ProductHandlerDeps{ProductRepository: productRepository},
	)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
		JWT:         jwtService,
	})

	// Middlewares
	stackMiddlewar := middleware.Chain(
		middleware.CORS,
		logger.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stackMiddlewar(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
