package model

type RegisterStore struct {
	StoreName   string `json:"store_name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type StoreResponse struct {
	UserID      uint   `json:"user_id"`
	StoreName   string `json:"store_name"`
	Description string `json:"description"`
}