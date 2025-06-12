package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lielamurs/balance-transactions/internal/dto"
	"github.com/lielamurs/balance-transactions/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

func (h *UserHandler) GetBalance(c echo.Context) error {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_user_id",
			Message: "User ID must be a positive integer",
		})
	}

	balance, err := h.userService.GetBalance(userID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "user_not_found",
				Message: "User does not exist",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to get balance",
			})
		}
	}

	return c.JSON(http.StatusOK, balance)
}

func (h *UserHandler) ProcessTransaction(c echo.Context) error {
	userID, sourceType, req, validationErr := h.validateTransactionRequest(c)
	if validationErr != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   validationErr.Code,
			Message: validationErr.Message,
		})
	}

	if err := h.userService.ProcessTransaction(userID, req, sourceType); err != nil {
		switch err.Error() {
		case "user not found":
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "user_not_found",
				Message: "User does not exist",
			})
		case "transaction already processed":
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "duplicate_transaction",
				Message: "Transaction with this ID has already been processed",
			})
		case "insufficient balance":
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "insufficient_balance",
				Message: "Account balance cannot be negative",
			})
		default:
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to process transaction",
			})
		}
	}

	return c.JSON(http.StatusOK, dto.TransactionResponse{
		Success: true,
		Message: "Transaction processed successfully",
	})
}

type ValidationError struct {
	Code    string
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (h *UserHandler) validateTransactionRequest(c echo.Context) (uint64, string, dto.TransactionRequest, *ValidationError) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "invalid_user_id",
			Message: "User ID must be a positive integer",
		}
	}

	sourceType := c.Request().Header.Get("Source-Type")
	if sourceType == "" {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "missing_header",
			Message: "Source-Type header is required",
		}
	}

	validSources := map[string]bool{"game": true, "server": true, "payment": true}
	if !validSources[sourceType] {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "invalid_source_type",
			Message: "Source-Type must be one of: game, server, payment",
		}
	}

	var req dto.TransactionRequest
	if err := c.Bind(&req); err != nil {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "invalid_request_body",
			Message: "Invalid JSON format",
		}
	}

	if req.State == "" {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "missing_state",
			Message: "State field is required",
		}
	}

	if req.State != "win" && req.State != "lose" {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "invalid_state",
			Message: "State must be 'win' or 'lose'",
		}
	}

	if req.TransactionID == "" {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "missing_transaction_id",
			Message: "TransactionId field is required",
		}
	}

	if err := h.validateAmount(req.Amount); err != nil {
		return 0, "", dto.TransactionRequest{}, &ValidationError{
			Code:    "invalid_amount",
			Message: err.Error(),
		}
	}

	return userID, sourceType, req, nil
}

func (h *UserHandler) validateAmount(amount string) error {
	if amount == "" {
		return errors.New("amount is required")
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return errors.New("invalid amount format")
	}

	if value <= 0 {
		return errors.New("amount must be positive")
	}

	parts := strings.Split(amount, ".")
	if len(parts) == 2 && len(parts[1]) > 2 {
		return errors.New("amount can have at most 2 decimal places")
	}

	return nil
}
