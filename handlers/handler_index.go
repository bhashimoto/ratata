package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func (cfg *ApiConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.getAccounts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error loading accounts")
		return
	}

	tmpFile := "./static/index.html"
	tmpl, err := template.ParseFiles(tmpFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to load template")
		return
	}
	data := struct {
		Accounts []Account
	}{
		Accounts: accounts,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
	}
	
}
