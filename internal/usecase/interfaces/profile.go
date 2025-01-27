package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type ProfileUseCase interface {
	CreateProfile(ctx context.Context, request *model.CreateProfile) (*model.ProfileResponse, error)
	GetProfile(ctx context.Context, id uint) (*model.ProfileResponse, error)
	UpdateProfile(ctx context.Context, request *model.UpdateProfile) (*model.ProfileResponse, error)
}