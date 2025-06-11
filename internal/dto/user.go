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
)
