package api

import (
	"context"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
	"github.com/bhashimoto/ratata/types"
)

func (cfg *ApiConfig) DBUserToUser(DBUser database.User) types.User {
	user := types.User{
		ID:         DBUser.ID,
		Name:       DBUser.Name,
		CreatedAt:  time.Unix(DBUser.CreatedAt, 0),
		ModifiedAt: time.Unix(DBUser.ModifiedAt, 0),
	}

	return user
}

func (cfg *ApiConfig) getUsersByAccount(accountID int64) ([]types.User, error) {
	dbUsers, err := cfg.DB.GetUsersByAccount(context.Background(), accountID)
	if err != nil {
		return []types.User{}, err
	}
	users := []types.User{}

	for _, dbuser := range dbUsers {
		user := cfg.DBUserToUser(dbuser)
		users = append(users, user)
	}
	return users, nil

}

func (cfg *ApiConfig) DBAccountToAccount(DBAccount database.Account) (types.Account, error) {
	transactions, err := cfg.getTransactionsByAccount(DBAccount.ID)
	if err != nil {
		return types.Account{}, err
	}

	users, err := cfg.getUsersByAccount(DBAccount.ID)
	if err != nil {
		return types.Account{}, err
	}

	account := types.Account{
		ID:           DBAccount.ID,
		Name:         DBAccount.Name,
		CreatedAt:    time.Unix(DBAccount.CreatedAt, 0),
		ModifiedAt:   time.Unix(DBAccount.ModifiedAt, 0),
		Transactions: transactions,
		Users:        users,
	}

	return account, nil
}

func (cfg *ApiConfig) DBTransactionToTransaction(dbt database.Transaction) (types.Transaction, error) {
	debts, err := cfg.getDebtsByTransaction(dbt.ID)
	if err != nil {
		return types.Transaction{}, err
	}

	dbUser, err := cfg.DB.GetUserByID(context.Background(), dbt.PaidBy)
	payer := cfg.DBUserToUser(dbUser)

	transaction := types.Transaction{
		ID:          dbt.ID,
		PaidBy:      payer,
		Description: dbt.Description,
		CreatedAt:   time.Unix(dbt.CreatedAt, 0),
		ModifiedAt:  time.Unix(dbt.ModifiedAt, 0),
		Amount:      dbt.Amount,
		Debts:       debts,
	}

	return transaction, nil
}


func (cfg *ApiConfig) DBDebtToDebt(dbd database.Debt) (types.Debt, error) {
	dbUser, err := cfg.DB.GetUserByID(context.Background(), dbd.UserID)
	if err != nil {
		return types.Debt{}, err
	}
	user := cfg.DBUserToUser(dbUser)

	debt := types.Debt{
		ID:            dbd.ID,
		User:          user,
		TransactionID: dbd.TransactionID,
		Amount:        dbd.Amount,
		CreatedAT:     time.Unix(dbd.CreatedAt, 0),
		ModifiedAt:    time.Unix(dbd.ModifiedAt, 0),
	}

	return debt, nil
}

func (cfg *ApiConfig) DBUserAccountToUserAccount(dbua database.UserAccount) (types.UserAccount, error) {
	userAccount := types.UserAccount{
		UserID:     dbua.UserID,
		AccountID:  dbua.AccountID,
		CreatedAt:  time.Unix(dbua.CreatedAt, 0),
		ModifiedAt: time.Unix(dbua.ModifiedAt, 0),
	}

	return userAccount, nil

}
