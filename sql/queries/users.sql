-- name: CreateUser :one
INSERT INTO users (email, hashed_passwords)
VALUES (
	$1,
	$2
	)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: FindUserFromEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one 
UPDATE users 
SET email = $1, hashed_passwords = $2,updated_at = NOW() 
WHERE id = $3
RETURNING *;

-- name: UpgradeUser :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING *;

