package types

import (
	"time"
)

type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Account struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	CreatedAt    time.Time     `json:"created_at"`
	ModifiedAt   time.Time     `json:"modified_at"`
	Users        []User        `json:"users"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	ID          int64     `json:"id"`
	PaidBy      User      `json:"paid_by"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	Amount      float64   `json:"amount"`
	Debts       []Debt    `json:"debts"`
}

type Debt struct {
	ID            int64     `json:"id"`
	User          User      `json:"user"`
	TransactionID int64     `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	CreatedAT     time.Time `json:"created_at"`
	ModifiedAt    time.Time `json:"modified_at"`
}

type UserAccount struct {
	UserID     int64     `json:"user_id"`
	AccountID  int64     `json:"account_id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
