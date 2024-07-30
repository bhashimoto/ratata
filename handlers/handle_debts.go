package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

func (cfg *ApiConfig) HandleDebtCreate(w http.ResponseWriter, r *http.Request) {
	params := struct{
		UserID int64 `json:"user_id"`
		TransactionID int64 `json:"transaction_id"`
		Amount float64 `json:"amount"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	dbd, err := cfg.DB.CreateDebt(r.Context(), database.CreateDebtParams{
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
		UserID: params.UserID,
		TransactionID: params.TransactionID,
		Amount: params.Amount,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create debt")
		return
	}

	debt := DBDebtToDebt(dbd)
	respondWithJSON(w, http.StatusCreated, debt)
}
