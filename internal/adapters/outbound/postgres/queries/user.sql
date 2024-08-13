-- name: GetUserById :one
SELECT id, email, username, password_hash, password_changed_at, role, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: UserExists :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE id = $1
);

-- name: GetUserByEmail :one
SELECT id, email, username, password_hash, password_changed_at, role, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash
)
VALUES ($1, $2, $3) 
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
    AND
    updated_at = sqlc.arg(updated_at)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
