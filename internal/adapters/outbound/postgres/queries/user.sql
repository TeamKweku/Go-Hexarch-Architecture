-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash,
    role
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash),
    role = COALESCE(sqlc.narg(role), role),
    password_changed_at = CASE 
        WHEN sqlc.narg(password_hash) IS NOT NULL THEN NOW()
        ELSE password_changed_at
  END
WHERE 
    id = sqlc.arg(id)
RETURNING *;