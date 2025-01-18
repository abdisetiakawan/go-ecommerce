package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func AuthToResponse(user *entity.User) *model.AuthResponse {
	return &model.AuthResponse{
		ID:          user.ID,
		UserUUID:    user.UserUUID,
		Username:    user.Username,
		Name:        user.Name,
		Email:       user.Email,
		Role:        user.Role,
		AccessToken: user.AccessToken,
		RefreshToken: user.RefreshToken,
	}
}