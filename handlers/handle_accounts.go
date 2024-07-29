package handlers

import (
	"encoding/json"
	"net/http"
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

	account := DBAccountToAccount(dbAccount)
	respondWithJSON(w, http.StatusCreated, account)
}

func (cfg *ApiConfig) HandleAccountsGet(w http.ResponseWriter, r *http.Request) {
	dbAccounts, err := cfg.DB.GetAccounts(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get accounts")
		return
	}

	accounts := []Account{}
	for _, dbAccount := range dbAccounts {
		accounts = append(accounts, DBAccountToAccount(dbAccount))
	}
	respondWithJSON(w, http.StatusOK, accounts)
}
