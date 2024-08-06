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

	tmpl := template.Must(template.ParseGlob("./static/*.html"))
	data := struct {
		Accounts []Account
	}{
		Accounts: accounts,
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
	}
	
}
