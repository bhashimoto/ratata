package routing

import (
	"net/http"

	"github.com/bhashimoto/ratata/handlers"
)

func SetRoutes(cfg *handlers.ApiConfig) (*http.ServeMux ){
	mux := http.NewServeMux()
	mux.HandleFunc("/", cfg.HandleIndex)

	mux.HandleFunc("GET /api/users", cfg.HandleGetUsers)
	mux.HandleFunc("GET /api/users/{userID}", cfg.HandleGetUser)
	mux.HandleFunc("POST /api/users", cfg.HandleCreateUser)

	mux.HandleFunc("POST /api/accounts", cfg.HandleAccountCreate)
	mux.HandleFunc("GET /api/accounts", cfg.HandleAccountsGet)
	mux.HandleFunc("GET /api/accounts/{accountID}", cfg.HandleAccountsGet)
	mux.HandleFunc("GET /api/accounts/{accountID}/balance", cfg.HandleBalanceGet)

	mux.HandleFunc("POST /api/transactions", cfg.HandleTransactionCreate)
	mux.HandleFunc("GET /api/transactions/{transactionID}", cfg.HandleTransactionGet)

	mux.HandleFunc("POST /api/debts", cfg.HandleDebtCreate)

	mux.HandleFunc("POST /api/user-accounts", cfg.HandleUserAccountCreate)

	return mux
}
