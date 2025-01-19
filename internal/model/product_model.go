package model

type RegisterProduct struct {
	AuthID      uint    `json:"-"`
	ProductName string  `json:"product_name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"required,min=10"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"required,gte=0"`
	Category    string  `json:"category" validate:"required,oneof=clothes electronics accessories"`
}

type ProductResponse struct {
	ProductUUID string  `json:"product_uuid"`
	StoreID     uint    `json:"store_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
}