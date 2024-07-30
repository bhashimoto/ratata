package handlers

import (
	"context"
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

	debt, err := cfg.DBDebtToDebt(dbd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve debt")
		return
	}

	respondWithJSON(w, http.StatusCreated, debt)
}

func (cfg *ApiConfig) getDebtsByTransaction(transactionID int64) ([]Debt, error) {
	dbDebts, err := cfg.DB.GetDebtsFromTransaction(context.Background(), transactionID)
	debts := []Debt{}
	if err != nil {
		return debts, err
	}

	for _, dbd := range dbDebts {
		debt, err := cfg.DBDebtToDebt(dbd)
		if err != nil {
			return []Debt{}, err
		}
		debts = append(debts, debt)
	}

	return debts, nil
}

