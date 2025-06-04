package product

import "github.com/lib/pq"

type ProductCreateRequest struct {
	Name        string         `json:"name"        validate:"required,min=3,max=20"`
	Description string         `json:"description" validate:"required"`
	Images      pq.StringArray `json:"images"      validate:"dive,base64"`
	Price       float64        `json:"price"       validate:"required,gt=0"`
}

type ProductUpdateRequest struct {
	Name        string         `json:"name"        validate:"min=3,max=20"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images"      validate:"dive,base64"`
	Price       float64        `json:"price"       validate:"gt=0"`
}

type PaggenationRequest struct {
	Page  int32 `json:"page"  validate:"number,gt=0"`
	Limit int32 `json:"limit" validate:"number,gt=0"`
}
