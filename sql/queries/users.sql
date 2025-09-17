-- name: CreateUser :one
INSERT INTO users (id, username, created_at, updated_at, hashed_password)
VALUES (?, ?, DATETIME('now'), DATETIME('now'), ?)
RETURNING id, username, created_at, updated_at;

-- name: GetUserByName :one
SELECT * FROM users
WHERE username = sqlc.arg(username);

-- name: SetUserDisplayName :one
UPDATE users
SET displayname = sqlc.arg(displayname), updated_at = DATETIME('now')
WHERE id = sqlc.arg(id)
RETURNING id, username, displayname, created_at, updated_at;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE id = ?;

