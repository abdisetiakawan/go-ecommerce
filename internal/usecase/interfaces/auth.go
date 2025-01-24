package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type AuthUseCase interface {
	Login(ctx context.Context, request *model.LoginUser) (*model.AuthResponse, error)
	Create(ctx context.Context, request *model.RegisterUser) (*model.AuthResponse, error)
}