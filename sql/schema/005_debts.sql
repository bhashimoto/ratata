-- +goose Up 
CREATE TABLE debts(
	id		INTEGER PRIMARY KEY AUTOINCREMENT,
	transaction_id	INTEGER NOT NULL,
	user_id		INTEGER NOT NULL,
	amount		REAL NOT NULL,
	created_at	INTEGER NOT NULL,
	modified_at	INTEGER NOT NULL,
	FOREIGN KEY(transaction_id) REFERENCES transactions(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE debts;
