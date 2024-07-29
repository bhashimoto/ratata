-- +goose Up
CREATE TABLE user_accounts(
	user_id INTEGER NOT NULL,
	account_id INTEGER NOT NULL,
		PRIMARY KEY(user_id, account_id)
);

-- +goose Down
DROP TABLE user_accounts;
