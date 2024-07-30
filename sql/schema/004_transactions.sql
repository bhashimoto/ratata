-- +goose Up
CREATE TABLE transactions(
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at	INTEGER NOT NULL,
	modified_at	INTEGER NOT NULL,
	description	TEXT NOT NULL,
	amount		REAL NOT NULL DEFAULT 0.0,
	paid_by		INTEGER NOT NULL,
	account_id	INTEGER NOT NULL,
	FOREIGN KEY(paid_by) REFERENCES users(id),
	FOREIGN KEY(account_id) REFERENCES accounts(id)
);

-- +goose Down
DROP TABLE transactions;
