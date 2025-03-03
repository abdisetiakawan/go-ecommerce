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
	StoreID     uint    `json:"store_id,omitempty"`
	Store       string  `json:"store_name,omitempty"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type GetProductsRequest struct {
	Role     string  `json:"-"`
	UserID   uint    `json:"-"`
	Search   string  `json:"-"`
	Category string  `json:"-" validate:"omitempty,oneof=clothes electronics accessories"`
	PriceMin float64 `json:"-" validate:"omitempty,gt=0"`
	PriceMax float64 `json:"-" validate:"omitempty,gt=0"`
	Page     int     `json:"-"`
	Limit    int     `json:"-"`
}

type GetProductRequest struct {
	UserID      uint   `json:"-" validate:"required"`
	ProductUUID string `json:"-" validate:"required,uuid"`
}

type UpdateProduct struct {
	UserID      uint    `json:"-" validate:"required"`
	ProductUUID string  `json:"-" validate:"required,uuid"`
	ProductName string  `json:"product_name" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,min=10"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	Stock       int     `json:"stock" validate:"omitempty,gte=0"`
	Category    string  `json:"category" validate:"omitempty,oneof=clothes electronics accessories"`
}

type DeleteProductRequest struct {
	UserID      uint   `json:"-" validate:"required"`
	ProductUUID string `json:"-" validate:"required,uuid"`
}