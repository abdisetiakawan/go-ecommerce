package model

type RegisterStore struct {
	ID          uint   `json:"-"`
	StoreName   string `json:"store_name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type StoreResponse struct {
	StoreName   string `json:"store_name"`
	Description string `json:"description"`
}

type UpdateStore struct {
	ID          uint   `json:"-"`
	StoreName   string `json:"store_name"`
	Description string `json:"description"`
}