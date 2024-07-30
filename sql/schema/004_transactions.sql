-- +goose Up
CREATE TABLE transactions(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created_at INTEGER NOT NULL,
	modified_at INTEGER NOT NULL,
	description TEXT NOT NULL,
	amount REAL DEFAULT 0.0,
	paid_by INTEGER NOT NULL,
	FOREIGN KEY(paid_by) REFERENCES users(id)
);

-- +goose Down
DROP TABLE transactions;
