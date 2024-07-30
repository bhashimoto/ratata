package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

func (cfg *ApiConfig) HandleTransactionCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Description	string  `json:"description"`
		Amount		float64 `json:"amount"`
		PaidBy		int64	`json:"paid_by"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	dbTransaction, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
		Amount: params.Amount,
		Description: params.Description,
		PaidBy: params.PaidBy,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating transaction")
		return
	}

	transaction := DBTransactionToTransaction(dbTransaction)
	respondWithJSON(w, http.StatusCreated, transaction)
}

func (cfg *ApiConfig) HandleTransactionGet(w http.ResponseWriter, r *http.Request) {
	transactionID, err  := strconv.Atoi(r.PathValue("transactionID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid transaction ID")
		return
	}

	dbt, err := cfg.DB.GetTransactionByID(r.Context(), int64(transactionID))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "transaction not found")
		return
	}

	dbDebts, err := cfg.DB.GetDebtsFromTransaction(r.Context(), int64(transactionID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error loading debts")
		return
	}


	transaction := DBTransactionToTransaction(dbt)
	debts := []Debt{}

	for _, dbd := range dbDebts {
		debts = append(debts, DBDebtToDebt(dbd))
	}

	ret := struct {
		Transaction Transaction `json:"transaction"`
		Debts []Debt `json:"debts"`
	}{
		Transaction: transaction,
		Debts: debts,
	}
	
	respondWithJSON(w, http.StatusOK, ret)
}
