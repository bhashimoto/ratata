package routing

import (
	"net/http"

	"github.com/bhashimoto/ratata/api"
	"github.com/bhashimoto/ratata/front"
)

func SetApiRoutes(cfg *api.ApiConfig, mux *http.ServeMux) error {
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// backend API endpoints
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

	return nil
}

func SetFrontEndRoutes(cfg *front.WebAppConfig, mux *http.ServeMux) error {
	//	fs := http.FileServer(http.Dir("./static"))
	mux.HandleFunc("/", cfg.HandleIndex)
	mux.HandleFunc("GET /accounts", cfg.HandleAccountsList)
	mux.HandleFunc("POST /accounts", cfg.HandleAccountCreate)
	mux.HandleFunc("/accounts/{accountID}", cfg.HandleAccounts)
	mux.HandleFunc("POST /accounts/{accountID}/create", cfg.HandleTransactionCreate)
	mux.HandleFunc("GET /accounts/{accountID}/create", cfg.HandleTransactionFormGet)
	mux.HandleFunc("GET /accounts/{accountID}/transactions", cfg.HandleTransactionsGet)
	//mux.HandleFunc("GET /accounts/{accountID}/payments", cfg.HandlePaymentsGet)
	return nil
}
