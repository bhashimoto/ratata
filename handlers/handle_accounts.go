package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

type Balance struct {
	User	User `json:"user"`
	Paid	float64 `json:"paid"`
	Owes	float64 `json:"owes"`
}

func (cfg *ApiConfig) HandleBalanceGet(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	transactions, err := cfg.getTransactionsByAccount(int64(accountID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve transactions")
		return
	}
	
	dbUsers, err := cfg.DB.GetUsersByAccount(r.Context(), int64(accountID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve users")
		return
	}

	balances := make(map[int64]*Balance)
	for _, dbUser := range dbUsers {
		user := cfg.DBUserToUser(dbUser)
		balance := Balance {
			User: user,
			Paid: 0.0,
			Owes: 0.0,
		}
		balances[user.ID] = &balance
	}

	for _, transaction := range transactions {
		if balances[transaction.PaidBy] == nil{
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("payer with id %v not registered in account", transaction.PaidBy))
			return
		}
		balances[transaction.PaidBy].Paid += transaction.Amount

		for _, debt := range transaction.Debts {
			if balances[debt.UserID] == nil {
				respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ower with id %v not registered in account",debt.UserID))
				return
			}
			balances[debt.UserID].Owes += debt.Amount
		}
	}

	respondWithJSON(w, http.StatusOK, balances)






}

func (cfg *ApiConfig) HandleAccountCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	dbAccount, err := cfg.DB.CreateAccount(r.Context(), database.CreateAccountParams{
		Name: params.Name,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create account")
		return
	}

	account, err := cfg.DBAccountToAccount(dbAccount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving account")
		return
	}
	respondWithJSON(w, http.StatusCreated, account)
}

func (cfg *ApiConfig) HandleAccountsGet(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.getAccounts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not load accounts")
		return
	}
	respondWithJSON(w, http.StatusOK, accounts)
}

func (cfg *ApiConfig) HandleAccountGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	account, err := cfg.getAccount(int64(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "account not found")
		return
	}

	dbUsers, err := cfg.DB.GetUsersByAccount(r.Context(), account.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to fetch account users")
		return
	}
	users := []User{}
	for _, dbUser := range dbUsers {
		user := cfg.DBUserToUser(dbUser)
		users = append(users, user)
	}

	resp := struct {
		Account Account `json:"account"`
		Users []User `json:"users"`
	}{
		Account: account,
		Users: users,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *ApiConfig) getAccounts() ([]Account, error) {
	dbAccounts, err := cfg.DB.GetAccounts(context.Background())
	if err != nil {
		return []Account{}, err
	}

	accounts := []Account{}
	for _, dbAccount := range dbAccounts {
		account, err := cfg.DBAccountToAccount(dbAccount)
		if err != nil {
			return []Account{}, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (cfg *ApiConfig) getAccount(id int64) (Account, error) {
	dbAccount, err := cfg.DB.GetAccountByID(context.Background(), id)
	if err != nil {
		return Account{}, err
	}

	account, err := cfg.DBAccountToAccount(dbAccount)
	if err != nil  {
		return Account{}, err
	}
	return account, nil
}
