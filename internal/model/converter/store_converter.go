package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func StoreToResponse(store *entity.Store) *model.StoreResponse {
	return &model.StoreResponse{
		StoreName:   store.StoreName,
		Description: store.Description,
		CreatedAt:   store.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   store.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}