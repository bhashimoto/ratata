package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

type AccountData struct {
	Account Account
	Balances map[User]*Balance
	Payments []Payment
}


func (cfg *ApiConfig) getAccountData(accountIDString string) (AccountData, error) {
	accountID, err := strconv.Atoi(accountIDString)
	if err != nil {
		return AccountData{}, err
	}
	
	account, err := cfg.getAccount(int64(accountID))
	if err != nil {
		return AccountData{}, err
	}

	balances, err := cfg.getBalancesFromAccount(int64(accountID))
	if err != nil {
		return AccountData{}, err
	}

	payments, err := cfg.calculatePayments(balances)
	if err != nil {
		return AccountData{}, err
	}

	ret := AccountData{
		Account: account,
		Balances: balances,
		Payments: payments,
	}

	return ret, nil

}

func (cfg *ApiConfig) FrontHandlePaymentsGet(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.getAccountData(r.PathValue("accountID"))
	tmpl, err := template.ParseFiles("./static/payments.html")
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tmpl.ExecuteTemplate(w, "payments", accountData)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	
}

func (cfg *ApiConfig) FrontHandleTransactionFormGet(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.getAccountData(r.PathValue("accountID"))
	tmpl, err := template.ParseFiles("./static/new_transaction.html")
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tmpl.Execute(w, accountData)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	
}

func (cfg *ApiConfig) FrontHandleTransactionCreate(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.getAccountData(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	paidBy, err := strconv.Atoi(r.Form.Get("payer"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	DBT, err := cfg.DB.CreateTransaction(r.Context(), database.CreateTransactionParams{
		AccountID: accountData.Account.ID,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
		Amount: amount,
		PaidBy: int64(paidBy),
		Description: r.Form.Get("description"),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	transaction, err := cfg.DBTransactionToTransaction(DBT)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payers := []int64{}
	for _, user := range accountData.Account.Users {
		if r.FormValue("checkbox-" + strconv.FormatInt(user.ID, 10)) == "on" {
			payers = append(payers, user.ID)
		}
	}
	params := []DebtParams{}
	for _, payer := range payers {
		params = append(params, DebtParams{
			UserId: payer,
			Amount: transaction.Amount / float64(len(payers)),
		})
	}
	_, err = cfg.insertDebts(transaction.ID, params)

	// Returning a new form and sending triggers
	tmpl, err := template.ParseFiles("./static/new_transaction.html")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accountData, _ = cfg.getAccountData(r.PathValue("accountID"))
	w.Header().Add("HX-Trigger", "newTransaction")
	err = tmpl.Execute(w, accountData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

func (cfg *ApiConfig) FrontHandleTransactionsGet(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.getAccountData(r.PathValue("accountID"))
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	tmpl, err := template.ParseFiles(
		"./static/transactions.html",
		"./static/new_transaction.html",
	)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = tmpl.ExecuteTemplate(w, "transactions" , accountData)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	
}

func (cfg *ApiConfig) FrontHandleAccounts(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.getAccountData(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	/*
	tmpl, err := template.ParseFiles(
		"./static/base.html",
		"./static/accounts.html",
		"./static/new_transaction.html",
		"./static/payments.html",
		"./static/transactions.html",
	)
	*/

	tmpl, err := template.ParseGlob("./static/*.html")

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}


	err = tmpl.ExecuteTemplate(w, "accounts", accountData)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

}
