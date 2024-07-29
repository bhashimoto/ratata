// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

type Account struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	ModifiedAt int64  `json:"modified_at"`
}

type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"created_at"`
	ModifiedAt int64  `json:"modified_at"`
}

type UserAccount struct {
	UserID    int64 `json:"user_id"`
	AccountID int64 `json:"account_id"`
}
