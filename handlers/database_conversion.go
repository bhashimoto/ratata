package handlers

import (
	"context"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

type User struct {
	ID	   int64     `json:"id"`
	Name	   string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func (cfg *ApiConfig) DBUserToUser(DBUser database.User) User {
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
	Users	     []User `json:"users"`
	Transactions []Transaction `json:"transactions"`
}

func (cfg *ApiConfig) getUsersByAccount(accountID int64) ([]User, error) {
	dbUsers, err := cfg.DB.GetUsersByAccount(context.Background(), accountID)
	if err != nil {
		return []User{}, err
	}
	users := []User{}

	for _, dbuser := range dbUsers {
		user := cfg.DBUserToUser(dbuser)
		users = append(users, user)
	}
	return users, nil

}

func (cfg *ApiConfig) DBAccountToAccount(DBAccount database.Account) (Account, error) {
	transactions, err := cfg.getTransactionsByAccount(DBAccount.ID)
	if err != nil {
		return Account{}, err
	}

	users, err := cfg.getUsersByAccount(DBAccount.ID)
	if err != nil {
		return Account{}, err
	}

	account := Account {
		ID: DBAccount.ID,
		Name: DBAccount.Name,
		CreatedAt: time.Unix(DBAccount.CreatedAt, 0),
		ModifiedAt: time.Unix(DBAccount.ModifiedAt, 0),
		Transactions: transactions,
		Users: users,
	}

	return account, nil
}

type Transaction struct {
	ID		int64		`json:"id"`
	PaidBy		User		`json:"paid_by"`
	Description	string		`json:"description"`
	CreatedAt	time.Time	`json:"created_at"`
	ModifiedAt	time.Time	`json:"modified_at"`
	Amount		float64		`json:"amount"`
	Debts		[]Debt		`json:"debts"`
}

func (cfg *ApiConfig) DBTransactionToTransaction(dbt database.Transaction) (Transaction, error) {
	debts, err := cfg.getDebtsByTransaction(dbt.ID)
	if err != nil {
		return Transaction{}, err
	}

	dbUser, err := cfg.DB.GetUserByID(context.Background(), dbt.PaidBy)
	payer := cfg.DBUserToUser(dbUser)

	transaction := Transaction {
		ID: dbt.ID,
		PaidBy: payer,
		Description: dbt.Description,
		CreatedAt: time.Unix(dbt.CreatedAt, 0),
		ModifiedAt: time.Unix(dbt.ModifiedAt, 0),
		Amount: dbt.Amount,
		Debts: debts, 
	}

	return transaction, nil
}

type Debt struct {
	ID int64 `json:"id"`
	User User`json:"user"`
	TransactionID int64 `json:"transaction_id"`
	Amount float64 `json:"amount"`
	CreatedAT time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func (cfg *ApiConfig) DBDebtToDebt(dbd database.Debt) (Debt, error) {
	dbUser, err := cfg.DB.GetUserByID(context.Background(), dbd.UserID)
	if err != nil {
		return Debt{}, err
	}
	user := cfg.DBUserToUser(dbUser)

	debt := Debt {
		ID: dbd.ID,
		User: user,
		TransactionID: dbd.TransactionID,
		Amount: dbd.Amount,
		CreatedAT: time.Unix(dbd.CreatedAt, 0),
		ModifiedAt: time.Unix(dbd.ModifiedAt, 0),
	}

	return debt, nil
}

type UserAccount struct {
	UserID int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`

}

func (cfg *ApiConfig) DBUserAccountToUserAccount(dbua database.UserAccount) (UserAccount, error) {
	userAccount := UserAccount {
		UserID: dbua.UserID,
		AccountID: dbua.AccountID,
		CreatedAt: time.Unix(dbua.CreatedAt, 0),
		ModifiedAt: time.Unix(dbua.ModifiedAt, 0),
	}

	return userAccount, nil

}

