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

	tmpl, err := template.ParseFiles(
		"./static/base.html", 
		"./static/index.html",
	)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	data := struct {
		Accounts []Account
	}{
		Accounts: accounts,
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	
}
