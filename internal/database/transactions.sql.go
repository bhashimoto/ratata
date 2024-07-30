// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: transactions.sql

package database

import (
	"context"
)

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (created_at, modified_at, description, amount, paid_by)
VALUES (?, ?, ?, ?, ?)
RETURNING id, created_at, modified_at, description, amount, paid_by, account_id
`

type CreateTransactionParams struct {
	CreatedAt   int64   `json:"created_at"`
	ModifiedAt  int64   `json:"modified_at"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	PaidBy      int64   `json:"paid_by"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.CreatedAt,
		arg.ModifiedAt,
		arg.Description,
		arg.Amount,
		arg.PaidBy,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Description,
		&i.Amount,
		&i.PaidBy,
		&i.AccountID,
	)
	return i, err
}

const getTransactionByID = `-- name: GetTransactionByID :one
SELECT id, created_at, modified_at, description, amount, paid_by, account_id
FROM transactions
WHERE id = ?
`

func (q *Queries) GetTransactionByID(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransactionByID, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.ModifiedAt,
		&i.Description,
		&i.Amount,
		&i.PaidBy,
		&i.AccountID,
	)
	return i, err
}

const getTransactionsByAccountID = `-- name: GetTransactionsByAccountID :many
SELECT id, created_at, modified_at, description, amount, paid_by, account_id
FROM transactions
WHERE account_id = ?
`

func (q *Queries) GetTransactionsByAccountID(ctx context.Context, accountID int64) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, getTransactionsByAccountID, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.ModifiedAt,
			&i.Description,
			&i.Amount,
			&i.PaidBy,
			&i.AccountID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
