package product

import (
	"4-order-api/pkg/db"
	"fmt"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	DataBase *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		DataBase: database,
	}
}

func (repo *ProductRepository) Create(product *Product) (*Product, error) {
	result := repo.DataBase.DB.Create(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (repo *ProductRepository) Update(product *Product) (*Product, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (repo *ProductRepository) Delete(id uint) error {
	result := repo.DataBase.DB.Delete(&Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("%s", "record not round")
	}

	return nil
}

func (repo *ProductRepository) GetById(id uint) (*Product, error) {
	var product Product
	result := repo.DataBase.DB.First(&product, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &product, nil
}

func (repo *ProductRepository) GetProducts(
	page uint,
	limit uint,
) ([]Product, error) {
	var products []Product
	result := repo.DataBase.DB.Offset(int((page - 1) * limit)).
		Limit(int(limit)).
		Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}
