package order

import "errors"

type OrderService struct {
	OrderRepository *OrderRepository
}

func NewOrderService(
	orderRepository *OrderRepository,
) *OrderService {
	return &OrderService{
		OrderRepository: orderRepository,
	}
}

func (service *OrderService) AddToCart(
	productId uint, userID uint,
) (uint, error) {
	product, _ := service.OrderRepository.GetProductById(productId)
	existedOrder, _ := service.OrderRepository.FindLastNotFormed(userID)
	if existedOrder == nil {
		order, err := service.OrderRepository.CreateNewOrder(product, userID)
		return order.ID, err
	} else {
		order, err := service.OrderRepository.AddProduct(existedOrder, product)
		return order.ID, err
	}
}

func (service *OrderService) FormedOrder(
	userID uint,
) (uint, error) {
	existedOrder, _ := service.OrderRepository.FindLastNotFormed(userID)
	if existedOrder == nil {
		return 0, errors.New(ErrorNotFindCartForOrder)
	}
	err := service.OrderRepository.OrderFormed(existedOrder)
	if err != nil {
		return 0, errors.New(ErrorCantFornedOrder)
	}
	return existedOrder.ID, nil
}
