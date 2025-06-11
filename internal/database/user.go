package database

import (
	"github.com/lielamurs/balance-transactions/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(userID uint64) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: GetDB(),
	}
}

func (r *userRepository) GetUser(userID uint64) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	return &user, err
}
