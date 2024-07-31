package routing

import (
	"net/http"

	"github.com/bhashimoto/ratata/handlers"
)

func SetRoutes(cfg *handlers.ApiConfig) (*http.ServeMux ){
	mux := http.NewServeMux()
	mux.HandleFunc("/", cfg.HandleIndex)

	mux.HandleFunc("GET /users", cfg.HandleGetUsers)
	mux.HandleFunc("GET /users/{userID}", cfg.HandleGetUser)
	mux.HandleFunc("POST /users", cfg.HandleCreateUser)

	mux.HandleFunc("POST /accounts", cfg.HandleAccountCreate)
	mux.HandleFunc("GET /accounts", cfg.HandleAccountsGet)
	mux.HandleFunc("GET /accounts/{accountID}", cfg.HandleAccountsGet)

	mux.HandleFunc("POST /transactions", cfg.HandleTransactionCreate)
	mux.HandleFunc("GET /transactions/{transactionID}", cfg.HandleTransactionGet)

	mux.HandleFunc("POST /debts", cfg.HandleDebtCreate)

	mux.HandleFunc("POST /user-accounts", cfg.HandleUserAccountCreate)

	return mux
}
