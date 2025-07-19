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
