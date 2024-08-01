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
	tmplPath := "./static/accounts.html"
	tmpl, _ := template.ParseFiles(tmplPath)

	data := struct{
		Account Account
	}{
		Account: account,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
	}

}
