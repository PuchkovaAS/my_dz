package product

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string
	Description string
	Images      pq.StringArray
	Price       float32
}
