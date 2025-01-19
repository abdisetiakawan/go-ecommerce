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

type AuthUseCase struct {
	DB *gorm.DB
	Log *logrus.Logger
	Validate *validator.Validate
	AuthRepository *repository.AuthRepository
	Jwt *helper.JwtHelper
	UUIDHelper *helper.UUIDHelper
}

func NewAuthUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, authRepos *repository.AuthRepository, jwt *helper.JwtHelper, uuid *helper.UUIDHelper) *AuthUseCase {
	return &AuthUseCase{
		DB: db,
		Log: log,
		Validate: validate,
		AuthRepository: authRepos,
		Jwt: jwt,
		UUIDHelper: uuid,
	}
}

func (c *AuthUseCase) Create(ctx context.Context, request *model.RegisterUser) (*model.AuthResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Failed to validate request body: %+v", err)
		return nil, err
	}

	// validate if email or username already exists
	total, err := c.AuthRepository.CountByField(tx, "email", request.Email)
	if err != nil {
		c.Log.Warnf("Failed to check email: %+v", err)
		return nil, err
	}
	if total > 0 {
		c.Log.Warnf("Email already exists")
		return nil, model.ErrUserAlreadyExists
	}
	ttl, err := c.AuthRepository.CountByField(tx, "username", request.Username)
	if err != nil {
		c.Log.Warnf("Failed to check username: %+v", err)
		return nil, err
	}
	if ttl > 0 {
		c.Log.Warnf("Username already exists")
		return nil, model.ErrUsernameExists
	}
	if request.Password != request.ConfirmPassword {
		c.Log.Warnf("Password and confirm password do not match")
		return nil, model.ErrPasswordNotMatch
	}
	// hash password
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to hash password: %+v", err)
		return nil, model.ErrInternalServer
	}

	user := &entity.User{
		UserUUID: c.UUIDHelper.Value,
		Username: request.Username,
		Name: request.Name,
		Email: request.Email,
		Role: request.Role,
		Password: string(password),
	}

	if err := c.AuthRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to create user: %+v", err)
		return nil, err
	}
	
	// generate token
	accessToken, refreshToken, err := c.Jwt.GenerateTokenUser(model.AuthResponse{
		ID: user.ID,
		Name: request.Name,
		Username: request.Username,
		Role: request.Role,
		Email: request.Email,
	})
	if err != nil {
		c.Log.Warnf("Failed to generate token: %+v", err)
		return nil, model.ErrInternalServer
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, err
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return converter.AuthToResponse(user), nil
}

func (c *AuthUseCase) Login(ctx context.Context, request *model.LoginUser) (*model.AuthResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Failed to validate request body: %+v", err)
		return nil, model.ErrBadRequest
	}

	// validate email
	user := new(entity.User)
    err := c.AuthRepository.FindByEmail(tx, user, request.Email)
    if err != nil {
        c.Log.Warnf("Failed to find user : %+v", err)
        return nil, model.ErrInvalidCredentials
    }

	// validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed to compare password: %+v", err)
		return nil, model.ErrInvalidCredentials
	}

	// generate token
	accessToken, refreshToken, err := c.Jwt.GenerateTokenUser(model.AuthResponse{
		ID: user.ID,
		Name: user.Name,
		Username: user.Username,
		Role: user.Role,
		Email: user.Email,
	})
	if err != nil {
		c.Log.Warnf("Failed to generate token: %+v", err)
		return nil, model.ErrInternalServer
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	err = tx.Commit().Error
	if err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, err
	}

	return converter.AuthToResponse(user), nil
}

