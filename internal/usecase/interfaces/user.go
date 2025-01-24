package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type UserUseCase interface {
	Create(ctx context.Context, request *model.CreateProfile) (*model.ProfileResponse, error)
	Get(ctx context.Context, userID uint) (*model.ProfileResponse, error)
	Update(ctx context.Context, request *model.UpdateProfile) (*model.ProfileResponse, error)
	ChangePassword(ctx context.Context, request *model.ChangePassword) error
}