package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type UserUseCase2 interface {
	Register(ctx context.Context, request *model.RegisterUser) (*model.AuthResponse, error)
	Login(ctx context.Context, request *model.LoginUser) (*model.AuthResponse, error)
	ChangePassword(ctx context.Context, request *model.ChangePassword) error
}