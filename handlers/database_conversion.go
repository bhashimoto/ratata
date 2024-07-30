package handlers

import (
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

type User struct {
	ID	   int64     `json:"id"`
	Name	   string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func DBUserToUser(DBUser database.User) User {
	user := User {
		ID: DBUser.ID,
		Name: DBUser.Name,
		CreatedAt: time.Unix(DBUser.CreatedAt, 0),
		ModifiedAt: time.Unix(DBUser.ModifiedAt, 0),
	}

	return user
}

type Account struct {
	ID	   int64     `json:"id"`
	Name	   string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func DBAccountToAccount(DBAccount database.Account) Account {
	account := Account {
		ID: DBAccount.ID,
		Name: DBAccount.Name,
		CreatedAt: time.Unix(DBAccount.CreatedAt, 0),
		ModifiedAt: time.Unix(DBAccount.ModifiedAt, 0),
	}

	return account
}

type Transaction struct {
	ID		int64		`json:"id"`
	PaidBy		int64		`json:"paid_by"`
	Description	string		`json:"description"`
	CreatedAt	time.Time	`json:"created_at"`
	ModifiedAt	time.Time	`json:"modified_at"`
	Amount		float64		`json:"amount"`
}

func DBTransactionToTransaction(dbt database.Transaction) Transaction {
	transaction := Transaction {
		ID: dbt.ID,
		PaidBy: dbt.PaidBy,
		Description: dbt.Description,
		CreatedAt: time.Unix(dbt.CreatedAt, 0),
		ModifiedAt: time.Unix(dbt.ModifiedAt, 0),
		Amount: dbt.Amount,
	}

	return transaction
}

type Debt struct {
	ID int64 `json:"id"`
	UserID int64 `json:"user_id"`
	TransactionID int64 `json:"transaction_id"`
	Amount float64 `json:"amount"`
	CreatedAT time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func DBDebtToDebt(dbd database.Debt) Debt {
	debt := Debt {
		ID: dbd.ID,
		UserID: dbd.UserID,
		TransactionID: dbd.TransactionID,
		Amount: dbd.Amount,
		CreatedAT: time.Unix(dbd.CreatedAt, 0),
		ModifiedAt: time.Unix(dbd.ModifiedAt, 0),
	}

	return debt
}
