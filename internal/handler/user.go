package handler

import (
	"net/http"
	"strconv"

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
