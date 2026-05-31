-- name: CreateUser :one
INSERT INTO users (id, email, display_name, password_hash, telegram_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET display_name = COALESCE(sqlc.narg('display_name'), display_name),
    telegram_id = CASE
        WHEN sqlc.narg('set_telegram')::boolean THEN sqlc.narg('telegram_id')
        ELSE telegram_id
    END
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT count(*) FROM users;