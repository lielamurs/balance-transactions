package model

import "time"

type (
	User struct {
		ID        uint64
		Balance   string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	Transaction struct {
		ID            uint64
		UserID        uint64
		TransactionID string
		Amount        string
		State         string
		SourceType    string
		CreatedAt     time.Time
	}
)
