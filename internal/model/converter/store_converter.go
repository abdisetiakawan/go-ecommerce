package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func StoreToResponse(store *entity.Store) *model.StoreResponse {
	return &model.StoreResponse{
		UserID:      store.UserID,
		StoreName:   store.StoreName,
		Description: store.Description,
	}
}