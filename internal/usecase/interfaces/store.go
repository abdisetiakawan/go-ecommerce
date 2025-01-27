package interfaces

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

type StoreUseCase interface {
	RegisterStore(ctx context.Context, request *model.RegisterStore) (*model.StoreResponse, error)
	GetStore(ctx context.Context, id uint) (*model.StoreResponse, error)
	UpdateStore(ctx context.Context, request *model.UpdateStore) (*model.StoreResponse, error)
}