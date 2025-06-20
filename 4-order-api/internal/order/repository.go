package order

import (
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type ProductWithQuantity struct {
	product.Product
	Quantity uint `json:"quantity"`
}

type OrderWithProducts struct {
	ID       uint                  `json:"id"`
	Products []ProductWithQuantity `json:"products"`
}

type OrderRepository struct {
	DataBase *db.Db
}

func NewOrderRepository(database *db.Db) *OrderRepository {
	return &OrderRepository{
		DataBase: database,
	}
}

func (repo *OrderRepository) GetProductById(id uint) (*product.Product, error) {
	var productObj product.Product
	result := repo.DataBase.DB.First(&productObj, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &productObj, nil
}

func (repo *OrderRepository) CreateNewOrder(
	productID uint, userID uint,
) (*Order, error) {
	order := Order{
		UserID:   userID,
		IsFormed: false,
	}

	err := repo.DataBase.DB.Create(&order).Error
	return &order, err
}

func (repo *OrderRepository) FindLastNotFormed(
	userID uint,
) (*Order, error) {
	var lastUserOrder Order

	result := repo.DataBase.DB.
		Where("user_id = ? AND is_formed = ?", userID, false).
		Order("created_at desc").
		First(&lastUserOrder).
		Error
	if errors.Is(result, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &lastUserOrder, result
}

func (repo *OrderRepository) AddProduct(
	order *Order, productID uint,
) error {
	var existingItem OrderProduct
	result := repo.DataBase.DB.
		Where("order_id = ? AND product_id = ?", order.ID, productID).
		First(&existingItem)

	if result.Error == nil {
		newQuantity := existingItem.Quantity + 1
		return repo.DataBase.DB.
			Model(&OrderProduct{}).
			Where("order_id = ? AND product_id = ?", order.ID, productID).
			Update("quantity", newQuantity).
			Error
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return repo.DataBase.DB.Create(&OrderProduct{
			OrderID:   order.ID,
			ProductID: productID,
			Quantity:  1,
		}).Error
	} else {
		return fmt.Errorf("database error: %w", result.Error)
	}
}

func (repo *OrderRepository) OrderFormed(
	order *Order,
) error {
	err := repo.DataBase.DB.Model(order).
		Update("is_formed", true).
		Error
	return err
}

func (repo *OrderRepository) GetAllOrders(
	userID uint,
) ([]uint, error) {
	var orders []uint

	err := repo.DataBase.DB.
		Table("orders").
		Select("orders.id").
		Where("user_id = ? AND is_formed = ?", userID, true).
		Scan(&orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (repo *OrderRepository) IsUsersOrder(
	orderID uint,
	userID uint,
) bool {
	var orderObj Order
	err := repo.DataBase.DB.
		Where(
			"id = ? AND user_id = ? AND is_formed = ?",
			orderID,
			userID,
			true,
		).
		First(&orderObj).
		Error
	if err != nil {
		return false
	}
	if orderObj.ID == 0 {
		return false
	}
	return true
}

func (repo *OrderRepository) GetOrderWithProducts(
	orderID uint,
	userID uint,
) (*OrderWithProducts, error) {
	type Result struct {
		product.Product
		Quantity uint
	}

	var results []Result

	err := repo.DataBase.DB.
		Table("order_products").
		Select("products.id, products.name, products.price, products.description, order_products.quantity").
		Joins("JOIN products ON products.id = order_products.product_id").
		Where("order_products.order_id = ?", orderID).
		Scan(&results).
		Error
	if err != nil {
		return nil, err
	}

	products := make([]ProductWithQuantity, len(results))
	for i, r := range results {
		products[i] = ProductWithQuantity{
			Product: product.Product{
				Name:        r.Name,
				Price:       r.Price,
				Description: r.Description,
			},
			Quantity: r.Quantity,
		}
	}

	var order Order
	if err := repo.DataBase.DB.First(&order, orderID).Error; err != nil {
		return nil, err
	}

	return &OrderWithProducts{
		ID:       order.ID,
		Products: products,
	}, nil
}
