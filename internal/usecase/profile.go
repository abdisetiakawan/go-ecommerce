package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	repo "github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProfileUseCase struct {
	db *gorm.DB
	val *validator.Validate
	profileRepo repo.ProfileRepository
}

func NewProfileUseCase(db *gorm.DB, val *validator.Validate, profileRepo repo.ProfileRepository) interfaces.ProfileUseCase {
	return &ProfileUseCase{
		db: db,
		val: val,
		profileRepo: profileRepo,
	}
}

func (u *ProfileUseCase) CreateProfile(ctx context.Context, request *model.CreateProfile) (*model.ProfileResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}

	profile := &entity.Profile{
		UserID: request.UserID,
		Gender: request.Gender,
		PhoneNumber: request.PhoneNumber,
		Address: request.Address,
		Avatar: request.Avatar,
		Bio: request.Bio,
	}
	hasProfile, err := u.profileRepo.CheckIDProfileByUserID(request.UserID)
	if err != nil {
		return nil, err
	}
	if hasProfile {
		return nil, model.ErrConflict
	}

	if err := u.profileRepo.CreateProfile(profile); err != nil {
		return nil, err
	}

	return converter.ProfileToResponse(profile), nil
}

func (u *ProfileUseCase) GetProfile(ctx context.Context, userID uint) (*model.ProfileResponse, error) {
	response, err := u.profileRepo.GetProfileByUserID(userID);
	if err != nil {
		return nil, err
	}
	return converter.ProfileToResponse(response), nil
}

func (u *ProfileUseCase) UpdateProfile(ctx context.Context, request *model.UpdateProfile) (*model.ProfileResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}

	existingProfile, err := u.profileRepo.GetProfileByUserID(request.UserID)
	if err != nil {
		return nil, err
	}

	if request.Gender != "" {
		existingProfile.Gender = request.Gender
	}
	if request.PhoneNumber != "" {
		existingProfile.PhoneNumber = request.PhoneNumber
	}
	if request.Address != "" {
		existingProfile.Address = request.Address
	}
	if request.Avatar != "" {
		existingProfile.Avatar = request.Avatar
	}
	if request.Bio != "" {
		existingProfile.Bio = request.Bio
	}

	if err := u.profileRepo.UpdateProfile(existingProfile); err != nil {
		return nil, err
	}

	return converter.ProfileToResponse(existingProfile), nil
}
