-- +goose up
CREATE TABLE users(
	id BLOB UNIQUE PRIMARY KEY,
	username TEXT NOT NULL,
	displayname TEXT,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose down
DROP TABLE users;
