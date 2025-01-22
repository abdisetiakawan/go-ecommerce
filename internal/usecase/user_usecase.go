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
	"golang.org/x/crypto/bcrypt"
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
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
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
	hasProfile, err := u.UserRepository.HasUserID(u.DB, request.UserID)
	if err != nil {
		u.Log.Warnf("Failed to check if user has profile: %+v", err)
		return nil, err
	}
	if hasProfile {
		u.Log.Warnf("User already has a profile")
		return nil, model.ErrConflict
	}

	if err := u.UserRepository.Create(u.DB, profile); err != nil {
		u.Log.Warnf("Failed to create profile: %+v", err)
		return nil, err
	}

	return converter.ProfileToResponse(profile), nil
}

func (u *UserUseCase) Get(ctx context.Context, userID uint) (*model.ProfileResponse, error) {
	var profile entity.Profile
	if err := u.UserRepository.FindByID(u.DB, &profile, userID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Profile not found")
			return nil, model.ErrNotFound
		}
		u.Log.Warnf("Failed to get profile: %+v", err)
		return nil, err
	}
	return converter.ProfileToResponse(&profile), nil
}

func (u *UserUseCase) Update(ctx context.Context, request *model.UpdateProfile) (*model.ProfileResponse, error) {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return nil, err
	}
	var profile entity.Profile
	if err := u.UserRepository.FindByID(u.DB, &profile, request.UserID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("Profile not found")
			return nil, model.ErrNotFound
		}
		u.Log.Warnf("Failed to get profile: %+v", err)
		return nil, err
	}
	if request.Gender != "" {
		profile.Gender = request.Gender
	}
	if request.PhoneNumber != "" {
		profile.PhoneNumber = request.PhoneNumber
	}
	if request.Address != "" {
		profile.Address = request.Address
	}
	if request.Avatar != "" {
		profile.Avatar = request.Avatar
	}
	if request.Bio != "" {
		profile.Bio = request.Bio
	}
	if err := u.UserRepository.Update(u.DB, &profile); err != nil {
		u.Log.Warnf("Failed to update profile: %+v", err)
		return nil, err
	}
	return converter.ProfileToResponse(&profile), nil
}

func (u *UserUseCase) ChangePassword(ctx context.Context, request *model.ChangePassword) error {
	if err := helper.ValidateStruct(u.Validate, u.Log, request); err != nil {
		return err
	}
	var user entity.User
	if err := u.UserRepository.AuthRepository.FindByID(u.DB, &user, request.UserID); err != nil {
		if err == gorm.ErrRecordNotFound {
			u.Log.Warnf("User not found")
			return model.ErrNotFound
		}
		u.Log.Warnf("Failed to get user: %+v", err)
		return err
	}
	if request.Password != request.ConfirmPassword {
		u.Log.Warnf("Password and confirm password do not match")
		return model.ErrPasswordNotMatch
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		u.Log.Warnf("Old password is incorrect")
		return model.ErrInvalidCredentials
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.Warnf("Failed to generate password: %+v", err)
		return err
	}
	user.Password = string(password)
	if err := u.UserRepository.AuthRepository.Update(u.DB, &user); err != nil {
		u.Log.Warnf("Failed to change password: %+v", err)
		return err
	}
	return nil
}