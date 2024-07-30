package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

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

	respondWithJSON(w, http.StatusOK, account)
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
