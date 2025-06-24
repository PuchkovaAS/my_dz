package main

import (
	"4-order-api/configs"
	"4-order-api/internal/auth"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"fmt"
	"net/http"
)

func App(conf *configs.Config) http.Handler {
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Logger
	logger := middleware.NewLogger(conf)

	// JWT
	jwtService := jwt.NewJWT(conf.Auth.Secret)

	// Repositorys
	productRepository := product.NewProductRepository(db)
	userRepository := user.NewUserRepository(db)
	orderRepository := order.NewOrderRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	orderService := order.NewOrderService(orderRepository)

	// Handlers
	product.NewProductHandler(
		router,
		product.ProductHandlerDeps{
			ProductRepository: productRepository,
			JWT:               jwtService,
			IOrderService:     orderService,
			IUserRepository:   userRepository,
		},
	)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
		JWT:         jwtService,
	})

	order.NewOrderHandler(
		router,
		order.OrderHandlerDeps{
			IUserRepository: userRepository,
			OrderService:    orderService,
			JWT:             jwtService,
		},
	)

	// Middlewares
	stackMiddlewar := middleware.Chain(
		middleware.CORS,
		logger.Logging,
	)

	return stackMiddlewar(router)
}

func main() {
	conf := configs.LoadConfig()
	app := App(conf)
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
