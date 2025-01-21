package model

type CreateProfile struct {
	UserID      uint   `json:"-"`
	Gender      string `json:"gender" validate:"required,oneof=male female"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Address     string `json:"address" validate:"required"`
	Avatar      string `json:"avatar" validate:"url,omitempty"`
	Bio         string `json:"bio"`
}

type ProfileResponse struct {
	UserID      uint   `json:"user_id"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Avatar      string `json:"avatar"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}