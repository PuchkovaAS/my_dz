package order

type OrderFormedResponse struct {
	OrderID uint `json:"orderID"`
}

type ProductResponse struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Quantity    uint    `json:"quantity"`
}

type OrderGetByIDResponse struct {
	OrderID  uint              `json:"orderID"`
	Products []ProductResponse `json:"products"`
}

type AllOrdersResponce struct {
	Orders []uint `json:"orders_id"`
}
