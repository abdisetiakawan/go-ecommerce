package model

type RegisterUser struct {
	Username        string `json:"username" validate:"required,min=3,max=100"`
	Name            string `json:"name" validate:"required,min=3,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Role            string `json:"role" validate:"required,oneof=seller buyer"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=seller buyer"`
}

type AuthResponse struct {
	ID           uint   `json:"id"`
	UserUUID     string `json:"user_uuid"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ChangePassword struct {
	UserID          uint   `json:"-" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8"`
	OldPassword     string `json:"old_password" validate:"required,min=8"`
}