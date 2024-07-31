package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)


func (cfg *ApiConfig) HandleUserAccountCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		UserID int64 `json:"user_id"`
		AccountID int64 `json:"account_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	dbUserAccount, err := cfg.DB.CreateUserAccount(r.Context(), database.CreateUserAccountParams{
		UserID: params.UserID,
		AccountID: params.AccountID,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user_account")
		return
	}

	userAccount, err := cfg.DBUserAccountToUserAccount(dbUserAccount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve user_account")
		return
	}
	respondWithJSON(w, http.StatusCreated, userAccount)
	

}
