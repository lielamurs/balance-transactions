package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/lielamurs/balance-transactions/internal/database"
	"github.com/lielamurs/balance-transactions/internal/dto"
	"github.com/lielamurs/balance-transactions/internal/model"
	"github.com/sirupsen/logrus"
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
	logrus.WithField("userID", userID).Info("Getting user balance")

	user, err := s.userRepo.GetUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithField("userID", userID).Warn("User not found")
			return nil, errors.New("user not found")
		}
		logrus.WithFields(logrus.Fields{"userID": userID, "error": err}).Error("Failed to get user")
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	balance, err := parseAmount(user.Balance)
	if err != nil {
		logrus.WithFields(logrus.Fields{"userID": userID, "balance": user.Balance, "error": err}).Error("Invalid balance format")
		return nil, fmt.Errorf("invalid balance format: %w", err)
	}

	formattedBalance := formatAmount(balance)
	logrus.WithFields(logrus.Fields{"userID": userID, "balance": formattedBalance}).Info("Balance retrieved successfully")
	return &dto.BalanceResponse{
		UserID:  user.ID,
		Balance: formattedBalance,
	}, nil
}

func (s *UserService) ProcessTransaction(userID uint64, req dto.TransactionRequest, sourceType string) error {
	logrus.WithFields(logrus.Fields{
		"userID":        userID,
		"transactionID": req.TransactionID,
		"state":         req.State,
		"amount":        req.Amount,
		"sourceType":    sourceType,
	}).Info("Starting transaction processing")

	return s.userRepo.GetDB().Transaction(func(tx *gorm.DB) error {
		exists, err := s.userRepo.TransactionExists(tx, req.TransactionID)
		if err != nil {
			return fmt.Errorf("failed to check existing transaction: %w", err)
		}
		if exists {
			logrus.WithFields(logrus.Fields{"userID": userID, "transactionID": req.TransactionID}).Warn("Duplicate transaction detected")
			return errors.New("transaction already processed")
		}

		user, err := s.userRepo.GetUserForUpdate(tx, userID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.WithFields(logrus.Fields{"userID": userID, "transactionID": req.TransactionID}).Warn("User not found for transaction")
				return errors.New("user not found")
			}
			logrus.WithFields(logrus.Fields{"userID": userID, "transactionID": req.TransactionID, "error": err}).Error("Failed to get user for transaction")
			return fmt.Errorf("failed to get user: %w", err)
		}

		currentBalance, err := parseAmount(user.Balance)
		if err != nil {
			return fmt.Errorf("invalid current balance: %w", err)
		}

		transactionAmount, err := parseAmount(req.Amount)
		if err != nil {
			return fmt.Errorf("invalid transaction amount: %w", err)
		}

		newBalance, err := calculateNewBalance(currentBalance, transactionAmount, req.State)
		if err != nil {
			if err.Error() == "insufficient balance" {
				logrus.WithFields(logrus.Fields{
					"userID":         userID,
					"transactionID":  req.TransactionID,
					"currentBalance": currentBalance,
					"amount":         transactionAmount,
				}).Warn("Insufficient balance for transaction")
			}
			return err
		}

		newBalanceStr := formatAmount(newBalance)
		if err := s.userRepo.UpdateUserBalance(tx, userID, newBalanceStr); err != nil {
			return fmt.Errorf("failed to update balance: %w", err)
		}

		transaction := &model.Transaction{
			UserID:        userID,
			TransactionID: req.TransactionID,
			Amount:        req.Amount,
			State:         req.State,
			SourceType:    sourceType,
		}

		if err := s.userRepo.CreateTransaction(tx, transaction); err != nil {
			logrus.WithFields(logrus.Fields{"userID": userID, "transactionID": req.TransactionID, "error": err}).Error("Failed to create transaction record")
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		logrus.WithFields(logrus.Fields{
			"userID":        userID,
			"transactionID": req.TransactionID,
			"oldBalance":    currentBalance,
			"newBalance":    newBalance,
		}).Info("Transaction processed successfully")
		return nil
	})
}

func parseAmount(amount string) (float64, error) {
	return strconv.ParseFloat(amount, 64)
}

func formatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}

func calculateNewBalance(currentBalance, transactionAmount float64, state string) (float64, error) {
	switch state {
	case "win":
		return currentBalance + transactionAmount, nil
	case "lose":
		newBalance := currentBalance - transactionAmount
		if newBalance < 0 {
			return 0, errors.New("insufficient balance")
		}
		return newBalance, nil
	default:
		return 0, errors.New("invalid transaction state")
	}
}
