package di

type IOrderService interface {
	AddToCart(productID, userID uint) (uint, error)
}

type IUserRepository interface {
	GetIdByPhone(phone string) (uint, error)
}
