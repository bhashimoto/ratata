package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
	"github.com/bhashimoto/ratata/types"
)

func (cfg *ApiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idString, err := strconv.Atoi(r.PathValue("userID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}
	dbUser, err := cfg.DB.GetUserByID(r.Context(), int64(idString))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	user := cfg.DBUserToUser(dbUser)

	dbAccs, err := cfg.DB.GetAccountsByUserID(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to retrieve accounts")
		return
	}
	accounts := []types.Account{}

	for _, dbAcc := range dbAccs {
		account, err := cfg.DBAccountToAccount(dbAcc)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "unable to get accounts")
			return
		}
		accounts = append(accounts, account)
	}

	resp := struct {
		User     types.User      `json:"user"`
		Accounts []types.Account `json:"accounts"`
	}{
		User:     user,
		Accounts: accounts,
	}

	respondWithJSON(w, http.StatusOK, resp)

}

func (cfg *ApiConfig) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch users")
		return
	}

	users := []types.User{}

	for _, dbUser := range dbUsers {
		users = append(users, cfg.DBUserToUser(dbUser))
	}

	respondWithJSON(w, http.StatusOK, users)
}

func (cfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:       params.Name,
		CreatedAt:  time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user")
		return
	}

	user := cfg.DBUserToUser(dbUser)

	respondWithJSON(w, http.StatusCreated, user)
}
