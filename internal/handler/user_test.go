package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lielamurs/balance-transactions/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestValidateAmount(t *testing.T) {
	handler := &UserHandler{}

	tests := []struct {
		name        string
		amount      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid decimal amount",
			amount:      "10.50",
			expectError: false,
		},
		{
			name:        "valid whole number",
			amount:      "100",
			expectError: false,
		},
		{
			name:        "valid small amount",
			amount:      "0.01",
			expectError: false,
		},
		{
			name:        "valid one decimal place",
			amount:      "5.5",
			expectError: false,
		},
		{
			name:        "valid two decimal places",
			amount:      "99.99",
			expectError: false,
		},
		{
			name:        "empty amount",
			amount:      "",
			expectError: true,
			errorMsg:    "amount is required",
		},
		{
			name:        "negative amount",
			amount:      "-10.50",
			expectError: true,
			errorMsg:    "amount must be positive",
		},
		{
			name:        "zero amount",
			amount:      "0",
			expectError: true,
			errorMsg:    "amount must be positive",
		},
		{
			name:        "zero decimal amount",
			amount:      "0.00",
			expectError: true,
			errorMsg:    "amount must be positive",
		},
		{
			name:        "invalid text",
			amount:      "abc",
			expectError: true,
			errorMsg:    "invalid amount format",
		},
		{
			name:        "mixed text and numbers",
			amount:      "10.50abc",
			expectError: true,
			errorMsg:    "invalid amount format",
		},
		{
			name:        "too many decimal places",
			amount:      "10.123",
			expectError: true,
			errorMsg:    "amount can have at most 2 decimal places",
		},
		{
			name:        "four decimal places",
			amount:      "5.1234",
			expectError: true,
			errorMsg:    "amount can have at most 2 decimal places",
		},
		{
			name:        "multiple dots",
			amount:      "10.50.25",
			expectError: true,
			errorMsg:    "invalid amount format",
		},
		{
			name:        "leading spaces",
			amount:      " 10.50",
			expectError: true,
			errorMsg:    "invalid amount format",
		},
		{
			name:        "trailing spaces",
			amount:      "10.50 ",
			expectError: true,
			errorMsg:    "invalid amount format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateAmount(tt.amount)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.EqualError(t, err, tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTransactionRequest(t *testing.T) {
	handler := &UserHandler{}

	tests := []struct {
		name           string
		userID         string
		sourceType     string
		jsonBody       string
		expectedUserID uint64
		expectedSource string
		expectedReq    dto.TransactionRequest
		expectedError  *ValidationError
	}{
		{
			name:           "valid request",
			userID:         "1",
			sourceType:     "game",
			jsonBody:       `{"state": "win", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedUserID: 1,
			expectedSource: "game",
			expectedReq: dto.TransactionRequest{
				State:         "win",
				Amount:        "10.50",
				TransactionID: "tx-001",
			},
			expectedError: nil,
		},
		{
			name:           "valid server source type",
			userID:         "2",
			sourceType:     "server",
			jsonBody:       `{"state": "lose", "amount": "5.00", "transactionId": "tx-002"}`,
			expectedUserID: 2,
			expectedSource: "server",
			expectedReq: dto.TransactionRequest{
				State:         "lose",
				Amount:        "5.00",
				TransactionID: "tx-002",
			},
			expectedError: nil,
		},
		{
			name:           "valid payment source type",
			userID:         "3",
			sourceType:     "payment",
			jsonBody:       `{"state": "win", "amount": "100.00", "transactionId": "tx-003"}`,
			expectedUserID: 3,
			expectedSource: "payment",
			expectedReq: dto.TransactionRequest{
				State:         "win",
				Amount:        "100.00",
				TransactionID: "tx-003",
			},
			expectedError: nil,
		},
		{
			name:       "invalid user ID - not a number",
			userID:     "abc",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_user_id",
				Message: "User ID must be a positive integer",
			},
		},
		{
			name:       "invalid user ID - negative",
			userID:     "-1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_user_id",
				Message: "User ID must be a positive integer",
			},
		},
		{
			name:       "missing source type header",
			userID:     "1",
			sourceType: "",
			jsonBody:   `{"state": "win", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "missing_header",
				Message: "Source-Type header is required",
			},
		},
		{
			name:       "invalid source type",
			userID:     "1",
			sourceType: "invalid",
			jsonBody:   `{"state": "win", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_source_type",
				Message: "Source-Type must be one of: game, server, payment",
			},
		},
		{
			name:       "invalid JSON",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"invalid": json}`,
			expectedError: &ValidationError{
				Code:    "invalid_request_body",
				Message: "Invalid JSON format",
			},
		},
		{
			name:       "missing state field",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "missing_state",
				Message: "State field is required",
			},
		},
		{
			name:       "empty state field",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "missing_state",
				Message: "State field is required",
			},
		},
		{
			name:       "invalid state value",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "invalid", "amount": "10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_state",
				Message: "State must be 'win' or 'lose'",
			},
		},
		{
			name:       "missing transaction ID",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "10.50"}`,
			expectedError: &ValidationError{
				Code:    "missing_transaction_id",
				Message: "TransactionId field is required",
			},
		},
		{
			name:       "empty transaction ID",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "10.50", "transactionId": ""}`,
			expectedError: &ValidationError{
				Code:    "missing_transaction_id",
				Message: "TransactionId field is required",
			},
		},
		{
			name:       "invalid amount - empty",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_amount",
				Message: "amount is required",
			},
		},
		{
			name:       "invalid amount - negative",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "-10.50", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_amount",
				Message: "amount must be positive",
			},
		},
		{
			name:       "invalid amount - too many decimals",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "10.123", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_amount",
				Message: "amount can have at most 2 decimal places",
			},
		},
		{
			name:       "invalid amount - not a number",
			userID:     "1",
			sourceType: "game",
			jsonBody:   `{"state": "win", "amount": "abc", "transactionId": "tx-001"}`,
			expectedError: &ValidationError{
				Code:    "invalid_amount",
				Message: "invalid amount format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Echo context
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/user/"+tt.userID+"/transaction", strings.NewReader(tt.jsonBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.sourceType != "" {
				req.Header.Set("Source-Type", tt.sourceType)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("userId")
			c.SetParamValues(tt.userID)

			userID, sourceType, request, validationErr := handler.validateTransactionRequest(c)

			if tt.expectedError != nil {
				assert.NotNil(t, validationErr)
				assert.Equal(t, tt.expectedError.Code, validationErr.Code)
				assert.Equal(t, tt.expectedError.Message, validationErr.Message)
				assert.Equal(t, uint64(0), userID)
				assert.Equal(t, "", sourceType)
			} else {
				assert.Nil(t, validationErr)
				assert.Equal(t, tt.expectedUserID, userID)
				assert.Equal(t, tt.expectedSource, sourceType)
				assert.Equal(t, tt.expectedReq.State, request.State)
				assert.Equal(t, tt.expectedReq.Amount, request.Amount)
				assert.Equal(t, tt.expectedReq.TransactionID, request.TransactionID)
			}
		})
	}
}
