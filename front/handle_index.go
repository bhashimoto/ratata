package front

import (
	"net/http"

	"github.com/bhashimoto/ratata/types"
)

func (cfg *WebAppConfig) HandleIndex(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.fetchAccounts()
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, "error loading accounts")
		return
	}

	//tmpl, err := template.ParseGlob("./static/*.html")
	data := struct {
		Title    string
		Accounts []types.Account
	}{
		Title:    "Welcome to Ratata",
		Accounts: accounts,
	}
	cfg.Templates.ExecuteTemplate(w, "index", data)
}
