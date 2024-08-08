package handlers

import (
	"bytes"
	"encoding/json"
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

func (cfg *ApiConfig) refreshAccountData(accountIdString string) (AccountData, error) {
	delete(cfg.AccountCache, accountIdString)
	return cfg.getAccountData(accountIdString)
}


func (cfg *ApiConfig) getAccountData(accountIDString string) (AccountData, error) {
	val, ok := cfg.AccountCache[accountIDString]
	if ok {
		return *val, nil
	}

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
	cfg.AddAccountDataToCache(accountIDString, &ret)

	return ret, nil

}

func (cfg *ApiConfig) FrontHandlePaymentsGet(w http.ResponseWriter, r *http.Request) {
	accountIDString := r.PathValue("accountID")
	accountData, err := cfg.getAccountData(accountIDString)
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
	accountIDString := r.PathValue("accountID")
	accountData, err := cfg.getAccountData(accountIDString)
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

	accountData, _ = cfg.refreshAccountData(accountIDString)
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

func (cfg *ApiConfig) FrontHandleAccountCreate(w http.ResponseWriter, r *http.Request) {
	// get form data
	// call backend endpoint
	// return success message with Hx-Trigger
	err := r.ParseForm()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := struct {
		Name string `json:"name"`
	}{
		Name: r.Form.Get("account_create_name"),
	}

	body, err := json.Marshal(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	bodyReader := bytes.NewReader(body)

	resp, err := http.Post("http://localhost:8080/api/accounts", "application/json", bodyReader)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.StatusCode >= 400 {
		log.Println("invalid status code")
		respondWithJSON(w, resp.StatusCode, "error")
		return
	}

	log.Println("adding header")
	w.Header().Add("Hx-Trigger", "newAccount")
	log.Println("serving accounts")
	cfg.serveAccounts(w)

}

func (cfg *ApiConfig) FrontHandleAccountsList(w http.ResponseWriter, r *http.Request) {
	cfg.serveAccounts(w)
}

func (cfg *ApiConfig) serveAccounts(w http.ResponseWriter) {
	acc, _ := cfg.getAccounts()
	accounts := struct {
		Accounts []Account
	}{
		Accounts: acc,
	}
	serveTempate(w, "./static/components/accounts_list.html", accounts)

}

func serveTempate(w http.ResponseWriter, filename string, data interface{}){
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		template.New("accounts").Execute(w, "error loading template")
		return
	}
	tmpl.Execute(w, data)
	
}
