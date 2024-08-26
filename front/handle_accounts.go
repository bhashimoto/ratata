package front

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/bhashimoto/ratata/types"
)

func (cfg *WebAppConfig) HandleTransactionFormGet(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.fetchAccount(r.PathValue("accountID"))
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cfg.Templates.ExecuteTemplate(w, "new_transaction", accountData)
}

func (cfg *WebAppConfig) HandleTransactionCreate(w http.ResponseWriter, r *http.Request) {
	accountIDString := r.PathValue("accountID")
	accountData, err := cfg.fetchAccount(accountIDString)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = r.ParseForm()
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	paidBy, err := strconv.Atoi(r.Form.Get("payer"))
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 64)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	params := struct {
		Description string       `json:"description"`
		Amount      float64      `json:"amount"`
		PaidBy      int64        `json:"paid_by"`
		AccountID   int64        `json:"account_id"`
	}{
		AccountID: accountData.ID,
		Description: r.Form.Get("description"),
		Amount: amount,
		PaidBy: int64(paidBy),
	}
	marshalled, err := json.Marshal(params)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	body := bytes.NewReader(marshalled)
	transactionResp, err := cfg.sendRequest("transactions", "POST", nil, body)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	transaction := types.Transaction{}
	decoder := json.NewDecoder(transactionResp.Body)
	err = decoder.Decode(&transaction)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	payers := []int64{}
	for _, user := range accountData.Users {
		if r.FormValue("checkbox-"+strconv.FormatInt(user.ID, 10)) == "on" {
			payers = append(payers, user.ID)
		}
	}
	type debtParams struct {
		UserID        int64   `json:"user_id"`
		TransactionID int64   `json:"transaction_id"`
		Amount        float64 `json:"amount"`
	}
	for _, payer := range payers {
		err = cfg.createDebt(payer, transaction.ID, (transaction.Amount / float64(len(payers))))
	}


	// Returning a new form and sending triggers
	tmpl, err := template.ParseFiles("./static/new_transaction.html")
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Add("HX-Trigger", "newTransaction")
	err = tmpl.Execute(w, accountData)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (cfg *WebAppConfig) createDebt(userId, transactionId int64, amount float64) error {
	type debtParams struct {
		UserID        int64   `json:"user_id"`
		TransactionID int64   `json:"transaction_id"`
		Amount        float64 `json:"amount"`
	}
	debt := debtParams{
		UserID: userId,
		TransactionID: transactionId,
		Amount: amount,
	}
	data, _ := json.Marshal(debt)
	body := bytes.NewReader(data)
	_ , err := cfg.sendRequest("debts", "POST", nil, body)
	if err != nil {
		return err
	}
	return nil
	
}

func (cfg *WebAppConfig) HandleTransactionsGet(w http.ResponseWriter, r *http.Request) {
	accountData, err := cfg.fetchAccount(r.PathValue("accountID"))
	if err != nil {
		log.Println(err)
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = cfg.Templates.ExecuteTemplate(w, "transactions", accountData)
	if err != nil {
		log.Println(err)
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (cfg *WebAppConfig) HandleAccounts(w http.ResponseWriter, r *http.Request) {
	accId := r.PathValue("accountID")
	accountData, err := cfg.fetchAccount(accId)
	if err != nil {
		log.Println("error from fetchAccount in HandleAccounts")
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	payments, balances, err := cfg.fetchAccountBalance(accId)
	data := struct {
		Account types.Account
		Balances []types.Balance
		Payments []types.Payment
	}{
		Account: accountData,
	}
	cfg.Templates.ExecuteTemplate(w, "accounts", data)



}

func (cfg *WebAppConfig) HandleAccountCreate(w http.ResponseWriter, r *http.Request) {
	// get form data
	// call backend endpoint
	// return success message with Hx-Trigger
	err := r.ParseForm()
	if err != nil {
		cfg.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	params := struct {
		Name string `json:"name"`
	}{
		Name: r.Form.Get("account_create_name"),
	}

	body, err := json.Marshal(params)
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	bodyReader := bytes.NewReader(body)

	resp, err := cfg.sendRequest("accounts", "POST", nil, bodyReader)
	defer resp.Body.Close()
	if err != nil {
		log.Println(err)
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if resp.StatusCode >= 400 {
		log.Println("invalid status code")
		cfg.RespondWithError(w, resp.StatusCode, "error")
		return
	}

	accounts, err := cfg.responseToAccounts(resp)
	if err != nil {
		log.Println("error at HandleAccountCreate:", err)
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Println("adding header")
	w.Header().Add("Hx-Trigger", "newAccount")
	log.Println("serving accounts")
	cfg.Templates.ExecuteTemplate(w, "accounts_list", accounts)
}

func (cfg *WebAppConfig) getAccount(id int64) {
	//TODO
}




func (cfg *WebAppConfig) HandleAccountsList(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.fetchAccounts()
	if err != nil {
		cfg.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	cfg.Templates.ExecuteTemplate(w, "components/accounts_list.html", accounts)
}

