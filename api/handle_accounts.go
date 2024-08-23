package api

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
	"github.com/bhashimoto/ratata/types"
)

type Balance struct {
	User types.User    `json:"user"`
	Paid float64 `json:"paid"`
	Owes float64 `json:"owes"`
}

type AccountData struct {
	Account  types.Account
	Balances map[types.User]*Balance
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
		Account:  account,
		Balances: balances,
		Payments: payments,
	}
	cfg.AddAccountDataToCache(accountIDString, &ret)

	return ret, nil

}

func (cfg *ApiConfig) getBalancesFromAccount(accountID int64) (map[types.User]*Balance, error) {
	transactions, err := cfg.getTransactionsByAccount(int64(accountID))
	if err != nil {
		return make(map[types.User]*Balance), err
	}

	dbUsers, err := cfg.DB.GetUsersByAccount(context.Background(), int64(accountID))
	if err != nil {
		return make(map[types.User]*Balance), err
	}

	balances := make(map[types.User]*Balance)
	for _, dbUser := range dbUsers {
		user := cfg.DBUserToUser(dbUser)
		balance := Balance{
			User: user,
			Paid: 0.0,
			Owes: 0.0,
		}
		balances[user] = &balance
	}

	for _, transaction := range transactions {
		if balances[transaction.PaidBy] == nil {
			return make(map[types.User]*Balance), err
		}
		balances[transaction.PaidBy].Paid += transaction.Amount

		for _, debt := range transaction.Debts {
			if balances[debt.User] == nil {
				return make(map[types.User]*Balance), err
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
		respondWithError(w, http.StatusInternalServerError, "unable to get balances: "+err.Error())
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
	From   types.User    `json:"from"`
	To     types.User    `json:"to"`
	Amount float64 `json:"amount"`
}

func (cfg *ApiConfig) calculatePayments(balances map[types.User]*Balance) ([]Payment, error) {
	if len(balances) == 0 {
		return []Payment{}, nil
	}

	type tally struct {
		user  types.User
		tally float64
	}

	tallies := []tally{}
	for _, balance := range balances {
		tallies = append(tallies, tally{
			user:  balance.User,
			tally: balance.Paid - balance.Owes,
		})
	}

	sort.Slice(tallies, func(i, j int) bool { return tallies[i].tally > tallies[j].tally })

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
				From:   tallies[from].user,
				To:     tallies[to].user,
				Amount: math.Abs(tallies[from].tally),
			})
			tallies[to].tally -= tallies[from].tally
			tallies[from].tally = 0
			from--
		} else if tallies[to].tally == tallies[from].tally {
			payments = append(payments, Payment{
				From:   tallies[from].user,
				To:     tallies[to].user,
				Amount: math.Abs(tallies[from].tally),
			})
			tallies[from].tally = 0
			tallies[to].tally = 0
			to++
			from--

		} else {
			payments = append(payments, Payment{
				From:   tallies[from].user,
				To:     tallies[to].user,
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
		Name:       params.Name,
		CreatedAt:  time.Now().Unix(),
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
	users := []types.User{}
	for _, dbUser := range dbUsers {
		user := cfg.DBUserToUser(dbUser)
		users = append(users, user)
	}

	resp := struct {
		Account types.Account `json:"account"`
		Users   []types.User  `json:"users"`
	}{
		Account: account,
		Users:   users,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *ApiConfig) getAccounts() ([]types.Account, error) {
	dbAccounts, err := cfg.DB.GetAccounts(context.Background())
	if err != nil {
		return []types.Account{}, err
	}

	accounts := []types.Account{}
	for _, dbAccount := range dbAccounts {
		account, err := cfg.DBAccountToAccount(dbAccount)
		if err != nil {
			return []types.Account{}, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (cfg *ApiConfig) getAccount(id int64) (types.Account, error) {
	dbAccount, err := cfg.DB.GetAccountByID(context.Background(), id)
	if err != nil {
		return types.Account{}, err
	}

	account, err := cfg.DBAccountToAccount(dbAccount)
	if err != nil {
		return types.Account{}, err
	}
	return account, nil
}
