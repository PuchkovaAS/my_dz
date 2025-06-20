package order

import (
	"4-order-api/pkg/di"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/response"
	"net/http"
	"strconv"
)

type OrderHandlerDeps struct {
	OrderService    *OrderService
	JWT             *jwt.JWT
	IUserRepository di.IUserRepository
}

type OrderHandler struct {
	OrderService    *OrderService
	IUserRepository di.IUserRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderService:    deps.OrderService,
		IUserRepository: deps.IUserRepository,
	}
	router.Handle(
		"POST /order",
		middleware.IsAuthed(handler.CreateOrder(), *deps.JWT),
	)

	router.Handle(
		"GET /order/{id}",
		middleware.IsAuthed(handler.GetOrderByID(), *deps.JWT),
	)

	router.Handle(
		"GET /my-orders",
		middleware.IsAuthed(handler.GetAllOrders(), *deps.JWT),
	)
}

func (handler *OrderHandler) CreateOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userPhone, ok := req.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, ErrorWrongJWT, http.StatusBadRequest)
			return
		}
		userID, err := handler.IUserRepository.GetIdByPhone(userPhone)
		if err != nil {
			http.Error(w, ErrorWrongJWT, http.StatusBadRequest)
			return
		}
		orderID, err := handler.OrderService.FormedOrder(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data := OrderFormedResponse{
			OrderID: orderID,
		}
		response.Json(w, data, http.StatusCreated)
	}
}

func (handler *OrderHandler) GetOrderByID() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userPhone, ok := req.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, ErrorWrongJWT, http.StatusBadRequest)
			return
		}
		userID, err := handler.IUserRepository.GetIdByPhone(userPhone)
		if err != nil {
			http.Error(w, ErrorWrongJWT, http.StatusBadRequest)
			return
		}

		idString := req.PathValue("id")
		orderID, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := handler.OrderService.OrderRepository.GetOrderWithProducts(
			uint(orderID),
			userID,
		)
		if err != nil {
			http.Error(
				w,
				"Internal server error",
				http.StatusInternalServerError,
			)
			return
		}
		data := &OrderGetByIDResponse{
			OrderID:  order.ID,
			Products: make([]ProductResponse, len(order.Products)),
		}

		for i, p := range order.Products {
			data.Products[i] = ProductResponse{
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			}
		}

		response.Json(w, data, http.StatusOK)
	}
}

func (handler *OrderHandler) GetAllOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userPhone, ok := req.Context().Value(middleware.ContextPhoneKey).(string)
		if !ok {
			http.Error(w, ErrorWrongJWT, http.StatusBadRequest)
			return
		}
		userID, err := handler.IUserRepository.GetIdByPhone(userPhone)
		if err != nil {
			http.Error(w, ErrorWrongJWT, http.StatusBadRequest)
			return
		}
		ordersId, err := handler.OrderService.OrderRepository.GetAllOrders(
			userID,
		)
		if err != nil {
			http.Error(
				w,
				"Internal server error",
				http.StatusInternalServerError,
			)
			return
		}

		var data AllOrdersResponce

		for _, orderId := range ordersId {
			data.Orders = append(data.Orders, orderId)
		}

		response.Json(w, data, http.StatusOK)
	}
}
