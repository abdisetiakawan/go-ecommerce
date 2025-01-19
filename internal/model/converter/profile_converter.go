package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func ProfileToResponse(profile *entity.Profile) *model.ProfileResponse {
	return &model.ProfileResponse{
		UserID:      profile.UserID,
		Gender:      profile.Gender,
		PhoneNumber: profile.PhoneNumber,
		Address:     profile.Address,
		Avatar:      profile.Avatar,
		Bio:         profile.Bio,
	}
}