package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func (cfg *ApiConfig) FrontHandleAccounts(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	account, err := cfg.getAccount(int64(accountID))
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, "error loading account")
		return
	}
	balances, err := cfg.getBalancesFromAccount(int64(accountID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	payments, err := cfg.calculatePayments(balances)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}


	tmpl, err := template.ParseFiles(
		"./static/base.html",
		"./static/accounts.html",
		"./static/new_transaction.html",
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}


	data := struct{
		Account Account
		Balances map[User]*Balance
		Payments []Payment
	}{
		Account: account,
		Balances: balances,
		Payments: payments,
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

}
