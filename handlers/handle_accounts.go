package handlers

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)


type Balance struct {
	User	User `json:"user"`
	Paid	float64 `json:"paid"`
	Owes	float64 `json:"owes"`
}

func (cfg *ApiConfig) getBalancesFromAccount(accountID int64) (map[User]*Balance, error)  {
	transactions, err := cfg.getTransactionsByAccount(int64(accountID))
	if err != nil {
		return make(map[User]*Balance), err
	}
	
	dbUsers, err := cfg.DB.GetUsersByAccount(context.Background(), int64(accountID))
	if err != nil {
		return make(map[User]*Balance), err
	}

	balances := make(map[User]*Balance)
	for _, dbUser := range dbUsers {
		user := cfg.DBUserToUser(dbUser)
		balance := Balance {
			User: user,
			Paid: 0.0,
			Owes: 0.0,
		}
		balances[user] = &balance
	}

	for _, transaction := range transactions {
		if balances[transaction.PaidBy] == nil{
			return make(map[User]*Balance), err
		}
		balances[transaction.PaidBy].Paid += transaction.Amount

		for _, debt := range transaction.Debts {
			if balances[debt.User] == nil {
				return make(map[User]*Balance), err
			}
			balances[debt.User].Owes += debt.Amount
		}
	}

	return balances, nil
	
}

func (cfg *ApiConfig) HandleBalanceGet(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.Atoi(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	
	balances, err := cfg.getBalancesFromAccount(int64(accountID))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to get balances: " + err.Error())
	}

	payments, err := cfg.calculatePayments(balances)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	ret_balance := []Balance{}
	for _, v := range balances {
		ret_balance = append(ret_balance, *v)
	}

	ret := struct {
		Balances []Balance `json:"balances"`
		Payments []Payment `json:"payments"`
	}{
		Balances: ret_balance,
		Payments: payments,
	}

	respondWithJSON(w, http.StatusOK, ret)
}

type Payment struct {
	From	User	`json:"from"`
	To	User	`json:"to"`
	Amount	float64 `json:"amount"`
}

func (cfg *ApiConfig) calculatePayments(balances map[User]*Balance) ([]Payment, error) {
	type tally struct {
		user User
		tally float64
	}

	tallies := []tally{}
	for _, balance := range balances {
		tallies = append(tallies, tally{
			user: balance.User,
			tally: balance.Paid - balance.Owes,
		})
	}

	sort.Slice(tallies, func(i, j int) bool {return tallies[i].tally > tallies[j].tally})

	payments := []Payment{}
	from := len(tallies) - 1
	to := 0
	for {
		if tallies[to].tally == 0.0 {
			break
		}
		if from == to {
			break
		}
		if tallies[from].tally == 0.0 {
			break
		}

		if tallies[to].tally > tallies[from].tally {
			payments = append(payments, Payment{
				From: tallies[from].user,
				To: tallies[to].user,
				Amount: math.Abs(tallies[from].tally),
			})
			tallies[to].tally -= tallies[from].tally
			tallies[from].tally = 0
			from--
		} else if tallies[to].tally == tallies[from].tally {
			payments = append(payments, Payment{
				From: tallies[from].user,
				To: tallies[to].user,
				Amount: math.Abs(tallies[from].tally),
			})
			tallies[from].tally = 0
			tallies[to].tally = 0
			to++
			from--
			
		} else {
			payments = append(payments, Payment{
				From: tallies[from].user,
				To: tallies[to].user,
				Amount: tallies[to].tally,
			})
			tallies[from].tally -= tallies[to].tally
			tallies[to].tally = 0
			to++

		}
	}
	return payments, nil
}

func (cfg *ApiConfig) HandleAccountCreate(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	dbAccount, err := cfg.DB.CreateAccount(r.Context(), database.CreateAccountParams{
		Name: params.Name,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create account")
		return
	}

	account, err := cfg.DBAccountToAccount(dbAccount)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving account")
		return
	}
	respondWithJSON(w, http.StatusCreated, account)
}

func (cfg *ApiConfig) HandleAccountsGet(w http.ResponseWriter, r *http.Request) {
	accounts, err := cfg.getAccounts()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not load accounts")
		return
	}
	respondWithJSON(w, http.StatusOK, accounts)
}

func (cfg *ApiConfig) HandleAccountGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("accountID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	account, err := cfg.getAccount(int64(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "account not found")
		return
	}

	dbUsers, err := cfg.DB.GetUsersByAccount(r.Context(), account.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to fetch account users")
		return
	}
	users := []User{}
	for _, dbUser := range dbUsers {
		user := cfg.DBUserToUser(dbUser)
		users = append(users, user)
	}

	resp := struct {
		Account Account `json:"account"`
		Users []User `json:"users"`
	}{
		Account: account,
		Users: users,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *ApiConfig) getAccounts() ([]Account, error) {
	dbAccounts, err := cfg.DB.GetAccounts(context.Background())
	if err != nil {
		return []Account{}, err
	}

	accounts := []Account{}
	for _, dbAccount := range dbAccounts {
		account, err := cfg.DBAccountToAccount(dbAccount)
		if err != nil {
			return []Account{}, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (cfg *ApiConfig) getAccount(id int64) (Account, error) {
	dbAccount, err := cfg.DB.GetAccountByID(context.Background(), id)
	if err != nil {
		return Account{}, err
	}

	account, err := cfg.DBAccountToAccount(dbAccount)
	if err != nil  {
		return Account{}, err
	}
	return account, nil
}
