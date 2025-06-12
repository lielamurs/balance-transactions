package database

import (
	"github.com/lielamurs/balance-transactions/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(userID uint64) (*model.User, error)
	GetUserForUpdate(tx *gorm.DB, userID uint64) (*model.User, error)
	UpdateUserBalance(tx *gorm.DB, userID uint64, newBalance string) error
	TransactionExists(tx *gorm.DB, transactionID string) (bool, error)
	CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error
	GetDB() *gorm.DB
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

func (r *userRepository) GetUserForUpdate(tx *gorm.DB, userID uint64) (*model.User, error) {
	var user model.User
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *userRepository) UpdateUserBalance(tx *gorm.DB, userID uint64, newBalance string) error {
	return tx.Model(&model.User{}).Where("id = ?", userID).Update("balance", newBalance).Error
}

func (r *userRepository) TransactionExists(tx *gorm.DB, transactionID string) (bool, error) {
	var count int64
	err := tx.Model(&model.Transaction{}).Where("transaction_id = ?", transactionID).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error {
	return tx.Create(transaction).Error
}

func (r *userRepository) GetDB() *gorm.DB {
	return r.db
}
