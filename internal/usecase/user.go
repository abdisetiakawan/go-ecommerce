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
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	db       *gorm.DB
	log      *logrus.Logger
	val      *validator.Validate
	userRepo repo.UserRepository
	uuid 	*helper.UUIDHelper
	jwt 	*helper.JwtHelper
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, userRepo repo.UserRepository, uuid *helper.UUIDHelper, jwt *helper.JwtHelper) interfaces.UserUseCase {
	return &UserUseCase{
		db:       db,
		log:      log,
		val:      validate,
		userRepo: userRepo,
		uuid: uuid,
		jwt: jwt,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, request *model.RegisterUser) (*model.AuthResponse, error) {
	if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
		return nil, err
	}

	// validate if email or username already exists
	total, err := uc.userRepo.IsUserFieldValueExist("email", request.Email)
	if err != nil {
		uc.log.Warnf("Failed to check email: %+v", err)
		return nil, err
	}
	if total {
		uc.log.Warnf("Email already exists")
		return nil, model.ErrUserAlreadyExists
	}
	total, err = uc.userRepo.IsUserFieldValueExist("username", request.Username)
	if err != nil {
		uc.log.Warnf("Failed to check username: %+v", err)
		return nil, err
	}
	if total {
		uc.log.Warnf("Username already exists")
		return nil, model.ErrUsernameExists
	}
	if request.Password != request.ConfirmPassword {
		uc.log.Warnf("Password and confirm password do not match")
		return nil, model.ErrPasswordNotMatch
	}
	// hash password
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Warnf("Failed to hash password: %+v", err)
		return nil, model.ErrInternalServer
	}

	user := &entity.User{
		UserUUID: uc.uuid.Generate(),
		Username: request.Username,
		Name: request.Name,
		Email: request.Email,
		Role: request.Role,
		Password: string(password),
	}

	if err := uc.userRepo.CreateUser(user); err != nil {
		uc.log.Warnf("Failed to create user: %+v", err)
		return nil, err
	}
	
	// generate token
	accessToken, refreshToken, err := uc.jwt.GenerateTokenUser(model.AuthResponse{
		ID: user.ID,
		Name: request.Name,
		Username: request.Username,
		Role: request.Role,
		Email: request.Email,
	})
	if err != nil {
		uc.log.Warnf("Failed to generate token: %+v", err)
		return nil, model.ErrInternalServer
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return converter.AuthToResponse(user), nil
}

func (uc *UserUseCase) Login(ctx context.Context, request *model.LoginUser) (*model.AuthResponse, error) {
	if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
		return nil, err
	}

    user, err := uc.userRepo.GetUserByEmail(request.Email)
    if err != nil {
        uc.log.Warnf("Failed to find user : %+v", err)
        return nil, model.ErrInvalidCredentials
    }

	// validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		uc.log.Warnf("Failed to compare password: %+v", err)
		return nil, model.ErrInvalidCredentials
	}

	// generate token
	accessToken, refreshToken, err := uc.jwt.GenerateTokenUser(model.AuthResponse{
		ID: user.ID,
		Name: user.Name,
		Username: user.Username,
		Role: user.Role,
		Email: user.Email,
	})
	if err != nil {
		uc.log.Warnf("Failed to generate token: %+v", err)
		return nil, model.ErrInternalServer
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return converter.AuthToResponse(user), nil
}



func (uc *UserUseCase) ChangePassword(ctx context.Context, request *model.ChangePassword) error {
	if err := helper.ValidateStruct(uc.val, uc.log, request); err != nil {
		return err
	}
	user, err := uc.userRepo.FindUserByID(request.UserID)
	if err != nil {
		uc.log.Warnf("Failed to find user: %+v", err)
		return err
	}
	if request.Password != request.ConfirmPassword {
		uc.log.Warnf("Password and confirm password do not match")
		return model.ErrPasswordNotMatch
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		uc.log.Warnf("Old password is incorrect")
		return model.ErrInvalidCredentials
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		uc.log.Warnf("Failed to generate password: %+v", err)
		return err
	}
	user.Password = string(password)
	if err := uc.userRepo.UpdateUser(user); err != nil {
		uc.log.Warnf("Failed to change password: %+v", err)
		return err
	}
	return nil
}

