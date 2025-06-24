package order

import (
	"errors"
)

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
	existedOrder, _ := service.OrderRepository.FindLastNotFormed(userID)
	if existedOrder == nil {
		existedOrder, err := service.OrderRepository.CreateNewOrder(
			productId,
			userID,
		)
		if err != nil {
			return 0, err
		}

		err = service.OrderRepository.AddProduct(existedOrder, productId)
		return existedOrder.ID, err

	} else {

		err := service.OrderRepository.AddProduct(existedOrder, productId)
		return existedOrder.ID, err
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
