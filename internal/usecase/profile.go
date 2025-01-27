package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProfileUseCase struct {
	db *gorm.DB
	log *logrus.Logger
	val *validator.Validate
	profileRepo *repository.ProfileRepository
}

func NewProfileUseCase(db *gorm.DB, log *logrus.Logger, val *validator.Validate, profileRepo *repository.ProfileRepository) interfaces.ProfileUseCase {
	return &ProfileUseCase{
		db: db,
		log: log,
		val: val,
		profileRepo: profileRepo,
	}
}

func (u *ProfileUseCase) CreateProfile(ctx context.Context, request *model.CreateProfile) (*model.ProfileResponse, error) {
	if err := helper.ValidateStruct(u.val, u.log, request); err != nil {
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
		u.log.Warnf("Failed to check if user has profile: %+v", err)
		return nil, err
	}
	if hasProfile {
		u.log.Warnf("User already has a profile")
		return nil, model.ErrConflict
	}

	if err := u.profileRepo.CreateProfile(profile); err != nil {
		u.log.Warnf("Failed to create profile: %+v", err)
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
	if err := helper.ValidateStruct(u.val, u.log, request); err != nil {
		return nil, err
	}
	var profile entity.Profile
	response, err := u.profileRepo.GetProfileByUserID(request.UserID)
	if err != nil {
		u.log.Warnf("Failed to get profile: %+v", err)
		return nil, err
	}
	if response.Gender != "" {
		profile.Gender = response.Gender
	}
	if response.PhoneNumber != "" {
		profile.PhoneNumber = response.PhoneNumber
	}
	if response.Address != "" {
		profile.Address = response.Address
	}
	if response.Avatar != "" {
		profile.Avatar = response.Avatar
	}
	if response.Bio != "" {
		profile.Bio = response.Bio
	}
	if err := u.profileRepo.UpdateProfile(&profile); err != nil {
		u.log.Warnf("Failed to update profile: %+v", err)
		return nil, err
	}
	return converter.ProfileToResponse(&profile), nil
}
