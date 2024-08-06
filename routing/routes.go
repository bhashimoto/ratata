package routing

import (
	"log"
	"net/http"

	"github.com/bhashimoto/ratata/handlers"
)

func SetRoutes(cfg *handlers.ApiConfig) (*http.ServeMux ){
	mux := http.NewServeMux()


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

	err := setFrontEndRoutes(cfg, mux)
	if err != nil {
		log.Fatal("Could not set frontend routes at SetRoutes")
	}
	return mux
}

func setFrontEndRoutes(cfg *handlers.ApiConfig, mux *http.ServeMux) (error) {
//	fs := http.FileServer(http.Dir("./static"))
	mux.HandleFunc("/", cfg.HandleIndex)
	mux.HandleFunc("/accounts/{accountID}", cfg.FrontHandleAccounts)
	mux.HandleFunc("POST /accounts/{accountID}/create", cfg.FrontHandleTransactionCreate)
	mux.HandleFunc("GET /accounts/{accountID}/create", cfg.FrontHandleTransactionFormGet)
	mux.HandleFunc("GET /accounts/{accountID}/transactions", cfg.FrontHandleTransactionsGet)
	mux.HandleFunc("GET /accounts/{accountID}/payments", cfg.FrontHandlePaymentsGet)
	return nil
}
