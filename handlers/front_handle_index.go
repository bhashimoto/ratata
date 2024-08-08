package handlers

import (
	//"html/template"
	"log"
	"net/http"
)

func (cfg *ApiConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.getAccounts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error loading accounts")
		return
	}

	//tmpl, err := template.ParseGlob("./static/*.html")
	tmpl, err := cfg.LoadTemplates("./static/")
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	data := struct {
		Title string
		Accounts []Account
	}{
		Title: "Welcome to Ratata",
		Accounts: accounts,
	}
	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	
}
