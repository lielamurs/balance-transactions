package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  string
		want    float64
		wantErr bool
	}{
		{
			name:    "valid decimal amount",
			amount:  "10.50",
			want:    10.50,
			wantErr: false,
		},
		{
			name:    "valid whole number",
			amount:  "100",
			want:    100.0,
			wantErr: false,
		},
		{
			name:    "valid small amount",
			amount:  "0.01",
			want:    0.01,
			wantErr: false,
		},
		{
			name:    "zero amount",
			amount:  "0",
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "zero decimal amount",
			amount:  "0.00",
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "invalid empty string",
			amount:  "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid text",
			amount:  "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid mixed text and numbers",
			amount:  "10.50abc",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseAmount(tt.amount)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestFormatAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		want   string
	}{
		{
			name:   "decimal amount",
			amount: 10.5,
			want:   "10.50",
		},
		{
			name:   "whole number",
			amount: 100,
			want:   "100.00",
		},
		{
			name:   "small amount",
			amount: 0.01,
			want:   "0.01",
		},
		{
			name:   "zero amount",
			amount: 0,
			want:   "0.00",
		},
		{
			name:   "large amount",
			amount: 999999.99,
			want:   "999999.99",
		},
		{
			name:   "amount with many decimals gets rounded",
			amount: 10.999,
			want:   "11.00",
		},
		{
			name:   "amount with rounding down",
			amount: 10.994,
			want:   "10.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAmount(tt.amount)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCalculateNewBalance(t *testing.T) {
	tests := []struct {
		name              string
		currentBalance    float64
		transactionAmount float64
		state             string
		wantBalance       float64
		wantErr           bool
		expectedError     string
	}{
		{
			name:              "win transaction adds amount",
			currentBalance:    100.0,
			transactionAmount: 50.0,
			state:             "win",
			wantBalance:       150.0,
			wantErr:           false,
		},
		{
			name:              "lose transaction subtracts amount",
			currentBalance:    100.0,
			transactionAmount: 30.0,
			state:             "lose",
			wantBalance:       70.0,
			wantErr:           false,
		},
		{
			name:              "lose transaction with exact balance",
			currentBalance:    50.0,
			transactionAmount: 50.0,
			state:             "lose",
			wantBalance:       0.0,
			wantErr:           false,
		},
		{
			name:              "lose transaction with insufficient balance",
			currentBalance:    30.0,
			transactionAmount: 50.0,
			state:             "lose",
			wantBalance:       0,
			wantErr:           true,
			expectedError:     "insufficient balance",
		},
		{
			name:              "win transaction with zero amount",
			currentBalance:    100.0,
			transactionAmount: 0.0,
			state:             "win",
			wantBalance:       100.0,
			wantErr:           false,
		},
		{
			name:              "lose transaction with zero amount",
			currentBalance:    100.0,
			transactionAmount: 0.0,
			state:             "lose",
			wantBalance:       100.0,
			wantErr:           false,
		},
		{
			name:              "win transaction from zero balance",
			currentBalance:    0.0,
			transactionAmount: 25.0,
			state:             "win",
			wantBalance:       25.0,
			wantErr:           false,
		},
		{
			name:              "lose transaction from zero balance",
			currentBalance:    0.0,
			transactionAmount: 10.0,
			state:             "lose",
			wantBalance:       0,
			wantErr:           true,
			expectedError:     "insufficient balance",
		},
		{
			name:              "invalid transaction state",
			currentBalance:    100.0,
			transactionAmount: 50.0,
			state:             "invalid",
			wantBalance:       0,
			wantErr:           true,
			expectedError:     "invalid transaction state",
		},
		{
			name:              "empty transaction state",
			currentBalance:    100.0,
			transactionAmount: 50.0,
			state:             "",
			wantBalance:       0,
			wantErr:           true,
			expectedError:     "invalid transaction state",
		},
		{
			name:              "decimal amounts work correctly",
			currentBalance:    10.50,
			transactionAmount: 5.25,
			state:             "win",
			wantBalance:       15.75,
			wantErr:           false,
		},
		{
			name:              "decimal lose transaction",
			currentBalance:    10.75,
			transactionAmount: 0.50,
			state:             "lose",
			wantBalance:       10.25,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculateNewBalance(tt.currentBalance, tt.transactionAmount, tt.state)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.EqualError(t, err, tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantBalance, got)
			}
		})
	}
}
