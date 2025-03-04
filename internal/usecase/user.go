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

	// Register registers a new user in the system and returns an authentication response.
	//
	// Parameters:
	//
	//   * ctx: context.Context - Context for the request, including the request body for user registration.
	//   * request: *model.RegisterUser - Request body for user registration.
	//
	// Returns:
	//
	//   * 201 Created: model.AuthResponse if user is registered successfully.
	//
	// Errors:
	//
	//   * 400 Bad Request: validation error or request body is invalid.
	//   * 409 Conflict: if email or username already exists.
	//   * 500 Internal Server Error: if there is an error when registering the user in the database.
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

// Login authenticates a user by validating their email and password.
// 
// It first validates the request structure. If the validation fails, it returns an error.
// Then, it retrieves the user by email. If the user is not found, it returns an error.
// It compares the hashed password stored for the user with the provided password.
// If the password is incorrect, it returns an error.
// If the password is correct, it generates a new access token and refresh token for the user.
// 
// Parameters:
// 
//   * ctx: context.Context - Context for the request.
//   * request: *model.LoginUser - Request body containing user login credentials.
// 
// Returns:
// 
//   * model.AuthResponse: Contains user data and tokens if authentication is successful.
//   * error: If an error occurs during the login process, such as validation failure or incorrect credentials.

func (uc *UserUseCase) Login(ctx context.Context, request *model.LoginUser) (*model.AuthResponse, error) {
	if err := helper.ValidateStruct(uc.val, request); err != nil {
		return nil, err
	}

    user, err := uc.userRepo.GetUserByEmail(request.Email)
    if err != nil {
        return nil, model.ErrInvalidCredentials
    }

	if user.Role != request.Role {
		return nil, model.ErrInvalidRole
	}	

	// validate password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
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



// ChangePassword changes a user's password if the old password is correct.
// It first checks if the new password and confirm password match.
// If they do not match, it returns an error.
// It then checks if the old password is correct.
// If the old password is incorrect, it returns an error.
// If the old password is correct, it updates the user's password and returns nil.
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

