package handlers

import (
	"time"

	"github.com/bhashimoto/ratata/internal/database"
)

type User struct {
	ID	   int64     `json:"id"`
	Name	   string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func DBUserToUser(DBUser database.User) User {
	user := User {
		ID: DBUser.ID,
		Name: DBUser.Name,
		CreatedAt: time.Unix(DBUser.CreatedAt, 0),
		ModifiedAt: time.Unix(DBUser.ModifiedAt, 0),
	}

	return user
}

type Account struct {
	ID	   int64     `json:"id"`
	Name	   string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func DBAccountToAccount(DBAccount database.Account) Account {
	account := Account {
		ID: DBAccount.ID,
		Name: DBAccount.Name,
		CreatedAt: time.Unix(DBAccount.CreatedAt, 0),
		ModifiedAt: time.Unix(DBAccount.ModifiedAt, 0),
	}

	return account
}
