package product

import (
	"4-order-api/pkg/di"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/request"
	"4-order-api/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
	JWT               *jwt.JWT
	IOrderService     di.IOrderService
	IUserRepository   di.IUserRepository
}

type ProductHandler struct {
	ProductRepository *ProductRepository
	IOrderService     di.IOrderService
	IUserRepository   di.IUserRepository
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
		IOrderService:     deps.IOrderService,
		IUserRepository:   deps.IUserRepository,
	}
	router.HandleFunc("POST /product", handler.Create())
	router.HandleFunc("PATCH /product/{id}", handler.Update())
	router.Handle(
		"POST /product/{id}/buy",
		middleware.IsAuthed(handler.Buy(), *deps.JWT),
	)
	router.HandleFunc("DELETE /product/{id}", handler.Delete())

	router.HandleFunc("GET /product/pagination", handler.Pagination())

	router.HandleFunc("GET /product/{id}", handler.Get())
}

func (handler *ProductHandler) Buy() http.HandlerFunc {
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
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		orderId, err := handler.IOrderService.AddToCart(uint(id), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data := ProductAddToCartResponse{
			ProductID: uint(id),
			OrderID:   orderId,
		}

		response.Json(w, data, http.StatusCreated)
	}
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[ProductCreateRequest](&w, req)
		if err != nil {
			return
		}
		product := &Product{
			Name:        body.Name,
			Description: body.Description,
			Price:       float32(body.Price),
			Images:      body.Images,
		}

		createdProduct, err := handler.ProductRepository.Create(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, createdProduct, http.StatusCreated)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[ProductUpdateRequest](&w, req)
		if err != nil {
			return
		}
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := handler.ProductRepository.Update(&Product{
			Model:       gorm.Model{ID: uint(id)},
			Name:        body.Name,
			Description: body.Description,
			Price:       float32(body.Price),
			Images:      body.Images,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, product, http.StatusCreated)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = handler.ProductRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, nil, http.StatusOK)
	}
}

func (handler *ProductHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product, err := handler.ProductRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, product, http.StatusOK)
	}
}

func (handler *ProductHandler) Pagination() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pageStr := req.URL.Query().Get("page")
		limitStr := req.URL.Query().Get("limit")

		if pageStr == "" {
			pageStr = "1"
		}
		if limitStr == "" {
			limitStr = "10"
		}
		page, err := strconv.ParseUint(pageStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		limit, err := strconv.ParseUint(limitStr, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		products, err := handler.ProductRepository.GetProducts(
			uint(page),
			uint(limit),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response.Json(w, products, http.StatusOK)
	}
}
