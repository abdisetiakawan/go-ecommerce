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

// CreateProfile creates a new user profile.
//
// It validates the request data and checks if a profile already exists for the given user ID.
// If a profile exists, it returns a conflict error. Otherwise, it creates a new profile
// and returns the created profile's response.
//
// Parameters:
//
//   * ctx: context.Context - Context for the request.
//   * request: *model.CreateProfile - Request body containing profile details.
//
// Returns:
//
//   * *model.ProfileResponse: Response body for the created profile.
//   * error: If validation fails, if a profile already exists, or if an error occurs during creation.

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

// GetProfile retrieves a profile by its user ID.
//
// It first checks if the profile exists. If the profile does not exist, it returns an error.
// If the profile exists, it returns the profile.
//
// Parameters:
//
//   * ctx: context.Context - Context for the request.
//   * userID: uint - User ID of the profile to be retrieved.
//
// Returns:
//
//   * *model.ProfileResponse: Response body for the profile.
//   * error: If an error occurs when retrieving a profile.
func (u *ProfileUseCase) GetProfile(ctx context.Context, userID uint) (*model.ProfileResponse, error) {
	response, err := u.profileRepo.GetProfileByUserID(userID);
	if err != nil {
		return nil, err
	}
	return converter.ProfileToResponse(response), nil
}

// UpdateProfile updates a profile by its user ID and the fields to be updated.
//
// It first checks if the profile exists. If the profile does not exist, it returns an error.
// If the profile exists, it updates the profile and returns the updated profile.
//
// Parameters:
//
//   * ctx: context.Context - Context for the request.
//   * request: *model.UpdateProfile - Request body for updating a profile.
//
// Returns:
//
//   * *model.ProfileResponse: Response body for the updated profile.
//   * error: If an error occurs when updating a profile.
func (u *ProfileUseCase) UpdateProfile(ctx context.Context, request *model.UpdateProfile) (*model.ProfileResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}

	existingProfile, err := u.profileRepo.GetProfileByUserID(request.UserID)
	if err != nil {
		return nil, err
	}

	if request.Name != "" {
		err := u.db.Model(&entity.User{}).Where("id = ?", request.UserID).Update("name", request.Name)
		if err.Error != nil {
			return nil, err.Error
		}
		existingProfile.Name = request.Name
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

	return converter.ProfileUpdatedToResponse(existingProfile), nil
}
