package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) Create(ctx context.Context, request *model.CreateProfile) (*model.ProfileResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

    if err := u.Validate.Struct(request); err != nil {
        if validationErrors, ok := err.(validator.ValidationErrors); ok {
            u.Log.Warnf("Validation failed: %+v", validationErrors)
            formattedErrors := helper.FormatValidationErrors(validationErrors)
            return nil, model.ErrValidationFailed(formattedErrors)
        }
        u.Log.Warnf("Failed to validate request body: %+v", err)
        return nil, model.ErrBadRequest
    }

	profile := &entity.Profile{
		UserID: request.UserID,
		Gender: request.Gender,
		PhoneNumber: request.PhoneNumber,
		Address: request.Address,
		Avatar: request.Avatar,
		Bio: request.Bio,
	}
	hasProfile, err := u.UserRepository.HasUserID(tx, request.UserID)
	if err != nil {
		u.Log.Warnf("Failed to check if user has profile: %+v", err)
		return nil, err
	}
	if hasProfile {
		u.Log.Warnf("User already has a profile")
		return nil, model.ErrConflict
	}

	if err := u.UserRepository.Create(tx, profile); err != nil {
		u.Log.Warnf("Failed to create profile: %+v", err)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, err
	}

	return converter.ProfileToResponse(profile), nil
}