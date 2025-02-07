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
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	db       *gorm.DB
	val      *validator.Validate
	userRepo repo.UserRepository
	uuid 	*helper.UUIDHelper
	jwt 	*helper.JwtHelper
}

func NewUserUseCase(db *gorm.DB, validate *validator.Validate, userRepo repo.UserRepository, uuid *helper.UUIDHelper, jwt *helper.JwtHelper) interfaces.UserUseCase {
	return &UserUseCase{
		db:       db,
		val:      validate,
		userRepo: userRepo,
		uuid: uuid,
		jwt: jwt,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, request *model.RegisterUser) (*model.AuthResponse, error) {
	if err := helper.ValidateStruct(uc.val, request); err != nil {
		return nil, err
	}

	// validate if email or username already exists
	total, err := uc.userRepo.IsUserFieldValueExist("email", request.Email)
	if err != nil {
		return nil, err
	}
	if total {
		return nil, model.ErrUserAlreadyExists
	}
	total, err = uc.userRepo.IsUserFieldValueExist("username", request.Username)
	if err != nil {
		return nil, err
	}
	if total {
		return nil, model.ErrUsernameExists
	}
	if request.Password != request.ConfirmPassword {
		return nil, model.ErrPasswordNotMatch
	}
	// hash password
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
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
		return nil, model.ErrInternalServer
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return converter.AuthToResponse(user), nil
}

func (uc *UserUseCase) Login(ctx context.Context, request *model.LoginUser) (*model.AuthResponse, error) {
	if err := helper.ValidateStruct(uc.val, request); err != nil {
		return nil, err
	}

    user, err := uc.userRepo.GetUserByEmail(request.Email)
    if err != nil {
        return nil, model.ErrInvalidCredentials
    }

	// validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
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
		return nil, model.ErrInternalServer
	}

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken

	return converter.AuthToResponse(user), nil
}



func (uc *UserUseCase) ChangePassword(ctx context.Context, request *model.ChangePassword) error {
	if err := helper.ValidateStruct(uc.val, request); err != nil {
		return err
	}
	user, err := uc.userRepo.FindUserByID(request.UserID)
	if err != nil {
		return err
	}
	if request.Password != request.ConfirmPassword {
		return model.ErrPasswordNotMatch
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword)); err != nil {
		return model.ErrInvalidCredentials
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(password)
	if err := uc.userRepo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}

