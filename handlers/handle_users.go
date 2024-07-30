package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/bhashimoto/ratata/internal/database"
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
	respondWithJSON(w, http.StatusOK, user)
	
}

func (cfg *ApiConfig) HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	dbUsers, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to fetch users")
		return
	}

	users := []User{}

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
		Name: params.Name,
		CreatedAt: time.Now().Unix(),
		ModifiedAt: time.Now().Unix(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user")
		return
	}

	user := cfg.DBUserToUser(dbUser)

	respondWithJSON(w, http.StatusCreated, user)
}
