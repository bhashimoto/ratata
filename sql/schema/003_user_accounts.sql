-- +goose Up
CREATE TABLE user_accounts(
	user_id INTEGER NOT NULL,
	account_id INTEGER NOT NULL,
	created_at INTEGER NOT NULL,
	modified_at INTEGER NOT NULL,
	PRIMARY KEY(user_id, account_id),
	FOREIGN KEY(user_id) REFERENCES users(id),
	FOREIGN KEY(account_id) REFERENCES accounts(id)
);

-- +goose Down
DROP TABLE user_accounts;
