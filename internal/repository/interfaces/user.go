package interfaces

import "github.com/abdisetiakawan/go-ecommerce/internal/entity"

type UserRepository interface {
	IsUserFieldValueExist(field string, value string) (bool, error)
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	FindUserByID(id uint) (*entity.User, error)
	UpdateUser(user *entity.User) error
}