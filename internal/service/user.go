package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lielamurs/balance-transactions/internal/database"
	"github.com/lielamurs/balance-transactions/internal/dto"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo database.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: database.NewUserRepository(),
	}
}

func (s *UserService) GetBalance(userID uint64) (*dto.BalanceResponse, error) {
	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	balance, err := strconv.ParseFloat(user.Balance, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid balance format: %w", err)
	}

	return &dto.BalanceResponse{
		UserID:  user.ID,
		Balance: fmt.Sprintf("%.2f", balance),
	}, nil
}
