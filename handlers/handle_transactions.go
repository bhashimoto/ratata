package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

func (cfg *ApiConfig) HandleTransactionCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Description	string		`json:"description"`
		Amount		float64		`json:"amount"`
		PaidBy		int64		`json:"paid_by"`
		AccountID	int64		`json:"account_id"`
		Debts		[]DebtParams	`json:"debts,omitempty"`
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
		AccountID: params.AccountID,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating transaction")
		return
	}

	_, err = cfg.insertDebts(dbTransaction.ID, params.Debts)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating debts")
		return
	}

	transaction, err := cfg.DBTransactionToTransaction(dbTransaction)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve transaction")
		return
	}



	respondWithJSON(w, http.StatusCreated, transaction)
}

func (cfg *ApiConfig) HandleTransactionGet(w http.ResponseWriter, r *http.Request) {
	transactionID, err  := strconv.Atoi(r.PathValue("transactionID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid transaction ID")
		return
	}

	transaction, err := cfg.getTransaction(int64(transactionID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error loading transaction")
		return
	}

	respondWithJSON(w, http.StatusOK, transaction)
}

func (cfg *ApiConfig) getTransaction(transactionID int64) (Transaction, error) {
	dbt, err := cfg.DB.GetTransactionByID(context.Background(), int64(transactionID))
	if err != nil {
		return Transaction{}, err
	}

	transaction, err := cfg.DBTransactionToTransaction(dbt)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (cfg *ApiConfig) getTransactionsByAccount(accountID int64) ([]Transaction, error) {
	dbts, err := cfg.DB.GetTransactionsByAccountID(context.Background(), accountID)
	if err != nil {
		return []Transaction{}, err
	}

	transactions := []Transaction{}
	for _, dbt := range dbts {
		transaction, err := cfg.DBTransactionToTransaction(dbt)
		if err != nil {
			return []Transaction{}, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
