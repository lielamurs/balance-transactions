package dto

type (
	BalanceResponse struct {
		UserID  uint64 `json:"userId"`
		Balance string `json:"balance"`
	}

	ErrorResponse struct {
		Error   string `json:"error"`
		Message string `json:"message,omitempty"`
	}

	TransactionRequest struct {
		State         string `json:"state" validate:"required,oneof=win lose"`
		Amount        string `json:"amount" validate:"required"`
		TransactionID string `json:"transactionId" validate:"required"`
	}

	TransactionResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
	}
)
