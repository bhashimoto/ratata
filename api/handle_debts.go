package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
	"github.com/bhashimoto/ratata/types"
)

type DebtParams struct {
	UserId int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

func (cfg *ApiConfig) HandleDebtCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		UserID        int64   `json:"user_id"`
		TransactionID int64   `json:"transaction_id"`
		Amount        float64 `json:"amount"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	dbd, err := cfg.DB.CreateDebt(r.Context(), database.CreateDebtParams{
		CreatedAt:     time.Now().Unix(),
		ModifiedAt:    time.Now().Unix(),
		UserID:        params.UserID,
		TransactionID: params.TransactionID,
		Amount:        params.Amount,
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

func (cfg *ApiConfig) getDebtsByTransaction(transactionID int64) ([]types.Debt, error) {
	dbDebts, err := cfg.DB.GetDebtsFromTransaction(context.Background(), transactionID)
	debts := []types.Debt{}
	if err != nil {
		return debts, err
	}

	for _, dbd := range dbDebts {
		debt, err := cfg.DBDebtToDebt(dbd)
		if err != nil {
			return []types.Debt{}, err
		}
		debts = append(debts, debt)
	}

	return debts, nil
}

func (cfg *ApiConfig) insertDebts(transactionID int64, params []DebtParams) ([]types.Debt, error) {
	debts := []types.Debt{}
	for _, protoDebt := range params {
		DBDebt, err := cfg.DB.CreateDebt(context.Background(), database.CreateDebtParams{
			UserID:        protoDebt.UserId,
			Amount:        protoDebt.Amount,
			CreatedAt:     time.Now().Unix(),
			ModifiedAt:    time.Now().Unix(),
			TransactionID: transactionID,
		})
		if err != nil {
			return []types.Debt{}, err
		}
		debt, err := cfg.DBDebtToDebt(DBDebt)
		if err != nil {
			return []types.Debt{}, err
		}
		debts = append(debts, debt)
	}
	return debts, nil
}
