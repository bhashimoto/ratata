-- +goose Up
CREATE TABLE users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	created_at INTEGER NOT NULL,
	modified_at INTEGER NOT NULL
);

-- +goose Down
DROP TABLE users;
