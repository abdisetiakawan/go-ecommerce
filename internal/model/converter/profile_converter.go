package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func ProfileToResponse(profile *entity.Profile) *model.ProfileResponse {
	return &model.ProfileResponse{
		UserID:      profile.UserID,
		Username:    profile.User.Username,
		Name:        profile.User.Name,
		Gender:      profile.Gender,
		PhoneNumber: profile.PhoneNumber,
		Address:     profile.Address,
		Avatar:      profile.Avatar,
		Bio:         profile.Bio,
		CreatedAt:   profile.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   profile.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ProfileUpdatedToResponse(profile *entity.Profile) *model.ProfileResponse {
	return &model.ProfileResponse{
		UserID:      profile.UserID,
		Username:    profile.User.Username,
		Name:        profile.Name,
		Gender:      profile.Gender,
		PhoneNumber: profile.PhoneNumber,
		Address:     profile.Address,
		Avatar:      profile.Avatar,
		Bio:         profile.Bio,
		CreatedAt:   profile.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   profile.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}